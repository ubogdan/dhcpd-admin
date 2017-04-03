package admin

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

var (
	serviceMngr      = "systemctl"
	fileConfig       = "/etc/dhcpd.conf"
	serviceStartArgs = []string{"start", "dhcpd4"}
	serviceStatus    = []string{"status", "dhcpd4"}
	serviceRestart   = []string{"restart", "dhcpd4"}
	serviceStop      = []string{"stop", "dhcpd4"}
)

type Service struct {
	lock sync.Mutex
}

//Start tells systemctl to start the dhcpd4 service
func (s *Service) Start() (string, error) {
	return s.run(serviceStartArgs...)
}

//Restart tells systemctl to restart the dhcpd4 service
func (s *Service) Restart() (string, error) {
	return s.run(serviceRestart...)
}

//Status tells systemctl to return the dhcpd4 service status
func (s *Service) Status() (string, error) {
	return s.run(serviceStatus...)
}

//Stop tells systemctl to stop the dhcpd4 service
func (s *Service) Stop() (string, error) {
	return s.run(serviceStop...)
}

func (s *Service) run(args ...string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cmd := exec.Command(serviceMngr, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

//UpdateConfig replaces the dhcpd4 config file for the new
//created by the config
func (s *Service) UpdateConfig(newConfig io.Reader) error {
	err := os.Rename(fileConfig, fileConfig+".backup")
	if err != nil {
		return err
	}
	f, err := os.Create(fileConfig)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, newConfig)
	return err
}
