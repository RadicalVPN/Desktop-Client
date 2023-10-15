//go:build windows

package platform

import (
	"path"
)

func initVariables() {
	serviceFile = path.Join("C:\\", "Program Files", "RadicalVPN", "service.txt")
	settingsFile = path.Join("C:\\", "Program Files", "RadicalVPN", "settings.json")
}
