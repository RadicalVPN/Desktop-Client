//go:build darwin

package platform

import (
	"path"
)

func initVariables() {
	settingsDirectory := "/Library/Application Support/RadicalVPN"

	serviceFile = path.Join(settingsDirectory, "service.txt")
	settingsFile = path.Join(settingsDirectory, "settings.json")
	wireguardPath = path.Join("")
	wireguardConfigPath = path.Join(settingsDirectory, "radicalvpn.conf")
}
