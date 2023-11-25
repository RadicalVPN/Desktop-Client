//go:build darwin

package wireguard

import (
	"fmt"
	"radicalvpnd/cli"
	"radicalvpnd/platform"
)

const (
	basePath = "/Applications/RadicalVPN.app/Contents/MacOS/Bash/bash"
)

func (wg *Wireguard) start() error {
	out, err := cli.Exec(basePath, platform.GetWireguardQuickPath(), "up", platform.GetWireguardConfPath())

	if err != nil {
		fmt.Println("err", err)
	}

	if out != "" {
		fmt.Println("out", out)
	}

	return nil
}

func (wg *Wireguard) stop() error {
	out, err := cli.Exec(basePath, platform.GetWireguardQuickPath(), "down", platform.GetWireguardConfPath())

	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println(out)

	return nil
}

func (wg *Wireguard) isConnected() bool {
	res, _ := cli.Exec(platform.GetWireguardPath(), "show")

	return res != ""
}
