package main

import (
	"fmt"
	"os"
	"os/exec"
	"printer/printer"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	SERVICE_NAME = "aaaPDFprint"
	PORT         = "8888"
)

type myService struct{}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	s <- svc.Status{State: svc.StartPending}
	p := &printer.Printer{}
	go p.RunServer(PORT)
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

func main() {
	if len(os.Args) > 1 && os.Args[1] == "i" {
		fmt.Println("Running in interactive mode")
		p := &printer.Printer{}
		go p.RunServer(PORT)
		for {
			time.Sleep(1 * time.Second)
		}
	} else {
		// Add firewall rule
		addFirewallRule()

		// Check if service is installed
		if !isServiceInstalled() {
			registerService()
		}

		runService()
	}
}

func runService() {
	el, err := eventlog.Open(SERVICE_NAME)
	if err != nil {
		return
	}
	defer el.Close()

	el.Info(1, fmt.Sprintf("%s service starting", SERVICE_NAME))
	err = svc.Run(SERVICE_NAME, &myService{})
	if err != nil {
		el.Error(1, fmt.Sprintf("%s service failed: %v", SERVICE_NAME, err))
		return
	}
	el.Info(1, fmt.Sprintf("%s service stopped", SERVICE_NAME))
}

func addFirewallRule() {
	// cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule", "name=\"Allow MyServiceName\"", "protocol=TCP", "dir=in", "localport=3000", "action=allow")
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=Allow %s", SERVICE_NAME),
		"protocol=TCP",
		"dir=in",
		fmt.Sprintf("localport=%s", PORT),
		"action=allow",
		"enable=yes")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to add firewall rule: %v\n", err)
	}
}

func isServiceInstalled() bool {
	cmd := exec.Command("sc", "query", SERVICE_NAME)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "RUNNING") || strings.Contains(string(output), "STOPPED")
}

func registerService() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get executable path: %v\n", err)
		return
	}

	cmd := exec.Command("sc", "create", SERVICE_NAME, "binPath=", fmt.Sprintf("\"%s\"", exePath), "start=auto")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to register service: %v\n", err)
		return
	}

	cmd = exec.Command("sc", "start", SERVICE_NAME)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to start service: %v\n", err)
	}
}
