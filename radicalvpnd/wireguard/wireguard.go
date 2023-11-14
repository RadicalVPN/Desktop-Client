package wireguard

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"radicalvpnd/cli"
	"radicalvpnd/logger"
	"radicalvpnd/platform"
	"radicalvpnd/settings"
	"radicalvpnd/util"
	"radicalvpnd/webapi"
	"sync"
)

type Wireguard struct {
}

var mutex sync.RWMutex

func NewWireguard() *Wireguard {
	return &Wireguard{}
}

var log *logger.Logger

func init() {
	log = logger.NewLogger("wireguard")

}

func (wg *Wireguard) downloadConfiguration(nodeLocation string, privacyFirewallLevel string) ([]byte, error) {
	dynamicvpn := webapi.DynamicVpnCreation{
		NodeLocation: nodeLocation,
		PrivacyFirewall: privacyFirewallLevel,
	}
	payload, _ := json.Marshal(dynamicvpn)

	req, err := http.NewRequest("PUT", "https://radicalvpn.com/api/1.0/internal/dynamic_vpn", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", settings.GetSessionCookie())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	//parse body
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	return body, nil
}

func (wg *Wireguard) Connect(nodeLocation string, privacyFirewallLevel string) error {
	mutex.Lock()
	defer mutex.Unlock()

	conf, err := wg.downloadConfiguration(nodeLocation, privacyFirewallLevel)
	if err != nil {
		return err
	}

	util.WriteFile(platform.GetWireguardConfPath(), conf, 0600)

	log.Info("Connecting to wireguard..")
	wg.start()

	return nil
}

func (wg *Wireguard) Disconnect() error {
	log.Info("Disconnecting from wireguard..")
	wg.stop()

	return nil
}

func (wg *Wireguard) IsConnected() bool {
	res, _ := cli.Exec(platform.GetWireguardPath(), "show")
	return res != ""
}
