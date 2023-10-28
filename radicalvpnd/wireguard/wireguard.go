package wireguard

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"radicalvpnd/logger"
	"radicalvpnd/platform"
	"radicalvpnd/settings"
	"radicalvpnd/util"
	"radicalvpnd/webapi"
)

type Wireguard struct {
}

func NewWireguard() *Wireguard {
	return &Wireguard{}
}

var log *logger.Logger

func init() {
	log = logger.NewLogger("wireguard")

}

func (wg *Wireguard) downloadConfiguration(node string) ([]byte, error) {
	dynamicvpn := webapi.DynamicVpnCreation{
		Node: "0102c185-9df6-4dbf-8580-272a68f67c69",
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

	return body, nil
}

func (wg *Wireguard) Connect(config string) error {
	conf, err := wg.downloadConfiguration("asd")
	if err != nil {
		return err
	}

	util.WriteFile(platform.GetWireguardConfPath(), conf, 0600)

	log.Info("Connecting to wireguard..")
	wg.start(config)

	return nil
}
