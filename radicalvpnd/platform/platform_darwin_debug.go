//go:build darwin && debug
// +build darwin,debug

package platform

import (
	"path"
)

func initVariables() {
	settingsDirectory := "/Library/Application Support/RadicalVPN"

	serviceFile = path.Join(settingsDirectory, "service.txt")
	settingsFile = path.Join(settingsDirectory, "settings.json")

	wireguardPath = path.Join("./deps/Darwin/Wireguard/wg")
	wireguardQuickPath = path.Join("./deps/Darwin/Wireguard/wg-quick.bash")
	wireguardConfigPath = path.Join(settingsDirectory, "radicalvpn.conf")

	logFilePath = path.Join(settingsDirectory, "radicalvpn.log")
}
