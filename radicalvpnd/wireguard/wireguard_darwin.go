//go:build darwin

package wireguard

import (
	"fmt"
	"radicalvpnd/cli"
	"radicalvpnd/platform"
)

func (wg *Wireguard) start(config string) error {
	out, err := cli.Exec(platform.GetWireguardQuickPath(), "up", platform.GetWireguardConfPath())

	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("out", out)

	return nil
}
