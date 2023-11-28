//go:build windows

package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log
var serviceName = "RadicalVPN-Daemon"

type radicalVpnService struct{}

func IsAdmin() bool {
	var sid *windows.SID

	//https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		log.Error(fmt.Sprintf("windows sid err: %s", err.Error()))
		return false
	}

	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		log.Error(fmt.Sprintf("token member error: %s", err.Error()))
		return false
	}

	return member
}

// this windows service stuff is inspired by https://github.com/golang/sys/blob/master/windows/svc/example/service.go
func PrepareRun() error {
	log.Info("Checking if daemon is started as a windows serivce..")

	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Error("failed to check if daemon is started as a windows service")
		return err
	}

	if !isService {
		log.Info("Daemon is not started as a windows service, starting as a regular process.. (DEVELOPOMENT ONLY)")
		return nil
	}

	log.Info("Starting daemon as a windows service..")

	go initWindowsService()

	return nil
}

func initWindowsService() {
	var err error

	elog, err = eventlog.Open(serviceName)
	if err != nil {
		log.Error(fmt.Sprintf("failed to open event log: %s", err.Error()))
		elog = nil
	}

	defer elog.Close()

	if elog != nil {
		elog.Info(1, fmt.Sprintf("starting %s service", serviceName))
	}

	err = svc.Run(serviceName, &radicalVpnService{})
	if err != nil {
		if elog != nil {
			elog.Error(1, fmt.Sprintf("%s service failed: %v", serviceName, err))
		}
		return
	}

	log.Info(fmt.Sprintf("%s service stopped", serviceName))
	if elog != nil {
		elog.Info(1, fmt.Sprintf("%s service stopped", serviceName))
	}
}

func (m *radicalVpnService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				testOutput := strings.Join(args, "-")
				testOutput += fmt.Sprintf("-%d", c.Context)
				elog.Info(1, testOutput)
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
