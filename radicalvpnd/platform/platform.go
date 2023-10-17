package platform

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	serviceFile   string
	settingsFile  string
	wireguardPath string
)

func init() {
	initVariables()
}

func Init() {

	//create service file
	if err := mkdir(filepath.Dir(serviceFile), os.ModePerm); err != nil {
		fmt.Println(err)
	}

	//create settings file
	if err := mkdir(filepath.Dir(settingsFile), os.ModePerm); err != nil {
		fmt.Println(err)
	}
}

func mkdir(path string, mode fs.FileMode) error {
	err := os.MkdirAll(path, mode)

	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	return nil
}

func GetServiceFIle() string {
	return serviceFile
}

func GetSettingsFile() string {
	return settingsFile
}

func GetWireguardPath() string {
	return wireguardPath
}
