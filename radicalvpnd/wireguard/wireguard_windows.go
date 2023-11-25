//go:build windows

package wireguard

import (
	"fmt"
	"path/filepath"
	"radicalvpnd/cli"
	"radicalvpnd/platform"
	"strings"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func (wg *Wireguard) getTunnelName() string {
	basePath := platform.GetWireguardConfPath()

	return strings.TrimSuffix(filepath.Base(basePath), filepath.Ext(basePath))
}

func (wg *Wireguard) getWireGuardServiceName() string {
	return "WireGuardTunnel$" + wg.getTunnelName()
}

func (wg *Wireguard) getServiceStatus(mgr *mgr.Mgr) (bool, svc.State, error) {
	service, err := mgr.OpenService(wg.getWireGuardServiceName())
	if err != nil {
		return false, 0, err
	}

	defer service.Close()

	state, err := service.Control(svc.Interrogate)
	if err != nil {
		return true, 0, err
	}

	return true, state.State, nil
}

func (wg *Wireguard) start() error {

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to service manager: %w", err)
	}

	fmt.Println("executing shit")

	cli.Exec(platform.GetWireguardPath(), "/installtunnelservice", platform.GetWireguardConfPath())

	//wait for service until installed
	isInstalled := false
	for {
		service, err := m.OpenService(wg.getWireGuardServiceName())
		if err == nil {
			fmt.Println("wireguard service successfully installed")
			isInstalled = true
			service.Close()
			break
		} else {
			fmt.Println("checking wireguard service status..")
		}
	}

	isStarted := false
	for {
		_, status, err := wg.getServiceStatus(m)
		if err != nil {
			return fmt.Errorf("service start failed")
		}

		if status == svc.Running {
			fmt.Println("Service is now running")
			isStarted = true
			break
		} else if status == svc.Stopped {
			fmt.Println("Service is stopped")
		}
	}

	fmt.Println(isInstalled, isStarted)

	defer m.Disconnect()

	return nil
}

func (wg *Wireguard) stop() error {
	return nil
}
