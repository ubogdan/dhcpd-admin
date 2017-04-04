package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/thewraven/dhcpd-admin"
	leases "github.com/thewraven/dhcpd-leases"
)

const leasesFile = "/var/lib/dhcp/dhcpd.leases"

var (
	globalTmpl, subnetTmpl, knownHostTmpl *template.Template
	logFile                               string
)

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func init() {
	var err error
	globalTmpl, err = template.New("global").Parse(admin.GlobalTmpl)
	checkError(err)
	knownHostTmpl, err = template.New("knownHost").Parse(admin.KnownHostTmpl)
	checkError(err)
	subnetTmpl, err = template.New("subnet").Parse(admin.SubnetTmpl)
	checkError(err)
}

func writeConfig(w io.Writer, global admin.GlobalConfig, subnets []admin.SubnetConfig,
	hosts []admin.HostConfig) {
	globalTmpl.Execute(w, global)
	for _, subnet := range subnets {
		subnetTmpl.Execute(w, subnet)
	}

	for _, host := range hosts {
		knownHostTmpl.Execute(w, host)
	}
}

type webServer struct {
	srv *admin.Service
}

func (ws *webServer) parseLeases(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	f, err := os.Open(leasesFile)
	if err != nil {
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	activeLeases, err := leases.ParseLeases(f)
	if err != nil {
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return

	}
	encoder.Encode(map[string]interface{}{
		"ok": activeLeases,
	})
}

func (ws *webServer) updateParams(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	config := &Config{}
	err := decoder.Decode(config)
	if err != nil {
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	result, err := updateAndRestart(config, ws.srv)
	if err != nil {
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	encoder.Encode(map[string]string{
		"ok": result,
	})
}

func (ws *webServer) getStatus(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	status, err := ws.srv.Status()
	if err != nil {
		encoder.Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	encoder.Encode(map[string]string{
		"ok": status,
	})
}

func (ws *webServer) watchLeases(every time.Duration, logFile string) {
	go func() {
		for {
			select {
			case <-time.After(every):
				fmt.Println("*****************	writing into log...")
				ws.dumpLeases(logFile)
			}

		}
	}()
}

func (ws webServer) dumpLeases(logFile string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 600)

	if err != nil {
		fmt.Println("cannot write log:", err)
		return
	}
	defer f.Close()
	/*	stats, err := f.Stat()
		if err != nil {
			fmt.Println("cannot read stats:", err)
			return
		}
		_, err = f.Seek(0, int(stats.Size()))
		if err != nil {
			fmt.Println("cannot seek:", err)
			return
		}*/
	encoder := json.NewEncoder(f)
	leaseFile, err := os.Open(leasesFile)
	if err != nil {
		fmt.Println("cannot read leases file:", err)
		return
	}
	defer leaseFile.Close()
	leases, err := leases.ParseLeases(leaseFile)
	if err != nil {
		fmt.Println("cannot parse leases", err)
		return
	}
	err = encoder.Encode(map[string]interface{}{
		"active": leases,
		"time":   time.Now(),
	})
	if err != nil {
		fmt.Println("cannot write in file:", err)
		return
	}
}

type Config struct {
	Global  admin.GlobalConfig
	Hosts   []admin.HostConfig
	Subnets []admin.SubnetConfig
}

func updateAndRestart(c *Config, srv *admin.Service) (string, error) {
	buffer := &bytes.Buffer{}
	writeConfig(buffer, c.Global, c.Subnets, c.Hosts)
	//	status, err := srv.Status()
	err := srv.UpdateConfig(buffer)
	if err != nil {
		return "", err
	}
	return srv.Restart()
}

func main() {
	/*global := admin.GlobalConfig{
		Lease: admin.Lease{
			Default: 60,
			Max:     120,
		},
		Authoritative: "authoritative",
		DNS1:          "8.8.8.8",
		DNS2:          "8.8.4.4",
	}
	hosts := []admin.HostConfig{{
		Host: "gerardo",
		MAC:  "1c:39:47:b0:6d:d1",
		IP:   "10.0.0.70",
	}}
	subnets := []admin.SubnetConfig{{
		Subnet:  "10.0.0.0",
		Netmask: "255.0.0.0",
		Known:   admin.IPRange{Initial: "10.0.0.1", Final: "10.0.0.5"},
		Unknown: admin.IPRange{Initial: "10.0.0.6", Final: "10.0.0.20"},
	}}
	buffer := &bytes.Buffer{}
	writeConfig(buffer, global, subnets, hosts)
	srv := &admin.Service{}
	status, err := srv.Status()
	fmt.Println(status, err)
	err = srv.UpdateConfig(buffer)*/
	mux := http.NewServeMux()
	ws := webServer{srv: &admin.Service{}}
	mux.HandleFunc("/status", ws.getStatus)
	mux.HandleFunc("/update", ws.updateParams)
	mux.HandleFunc("/leases", ws.parseLeases)
	ws.watchLeases(time.Second*10, "dhcpd.log")
	server := http.Server{}
	server.Addr = ":9000"
	server.Handler = mux
	fmt.Println("Listening at :9000...")
	server.ListenAndServe()
}
