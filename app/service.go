package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

type Service struct {
	Name string
	Port string
}

func (m *Service) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	s <- svc.Status{State: svc.StartPending}
	p := &Printer{}
	go p.RunServer(m.Port)
	s <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				s <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s <- svc.Status{State: svc.StopPending}
				break loop
			default:
				// unexpected control request, just ignore it
			}
		}
	}

	s <- svc.Status{State: svc.Stopped}
	return
}

func (m *Service) RunProd() {
	el, err := eventlog.Open(m.Name)
	if err != nil {
		return
	}
	defer el.Close()

	err = m.AddFirewallRule()
	if err != nil {
		el.Info(1, fmt.Sprintf("%s error: %v", m.Name, err))
	}

	if !m.IsInstalled() {
		err = m.Register()
		if err != nil {
			el.Info(1, fmt.Sprintf("%s error: %v", m.Name, err))
		}
	}

	el.Info(1, fmt.Sprintf("%s service starting", m.Name))
	err = svc.Run(m.Name, &Service{})
	if err != nil {
		el.Error(1, fmt.Sprintf("%s service failed: %v", m.Name, err))
		return
	}
	el.Info(1, fmt.Sprintf("%s service stopped", m.Name))
}

func (m *Service) RunDev() {
	fmt.Println("Running in dev mode")
	p := &Printer{}
	go p.RunServer(m.Port)
	for {
		time.Sleep(1 * time.Second)
	}
}

func (m *Service) AddFirewallRule() error {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=Allow %s", m.Name),
		"protocol=TCP",
		"dir=in",
		fmt.Sprintf("localport=%s", m.Port),
		"action=allow",
		"enable=yes")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add firewall rule: %v", err)
	}

	return nil
}

func (m *Service) IsInstalled() bool {
	cmd := exec.Command("sc", "query", m.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "RUNNING") || strings.Contains(string(output), "STOPPED")
}

func (m *Service) Register() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	cmd := exec.Command("sc", "create", m.Name, "binPath=", fmt.Sprintf("\"%s\"", exePath), "start=auto")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	cmd = exec.Command("sc", "start", m.Name)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	return nil
}
