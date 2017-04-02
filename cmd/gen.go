package main

import (
	"io"
	"os"
	"text/template"

	"github.com/thewraven/dhcpd-admin"
)

var (
	globalTmpl, subnetTmpl, knownHostTmpl *template.Template
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

func main() {
	global := admin.GlobalConfig{
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
	writeConfig(os.Stdout, global, subnets, hosts)
}
