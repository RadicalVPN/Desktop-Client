//go:build windows

package wireguard

import (
	"fmt"
	"path/filepath"
	"radicalvpnd/cli"
	"radicalvpnd/platform"
	"strings"
	"sync"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	connectionMutex sync.Mutex
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
	connectionMutex.Lock()
	defer func() {
		connectionMutex.Unlock()
	}()

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to service manager: %w", err)
	}

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
	connectionMutex.Lock()
	defer func() {
		connectionMutex.Unlock()
	}()

	m, err := mgr.Connect()
	if err != nil {
		log.Error("connection to service manager failed")
		return fmt.Errorf("could not connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(wg.getWireGuardServiceName())
	if err != nil {
		log.Warning("stop wireguard shutdown, service not available")
		return nil // service not available (so, nothing to uninstall)
	}
	s.Close()

	res, err := cli.Exec(platform.GetWireguardPath(), "/uninstalltunnelservice", wg.getTunnelName())

	log.Info("wireguard disconnected", res)

	return nil
}
