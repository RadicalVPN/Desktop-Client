package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"radicalvpnd/platform"
	"radicalvpnd/util"
	"sync"
)

type Settings struct {
	Session Session `json:"session"`
}

var mutx sync.RWMutex

func NewSettings() *Settings {
	settings := &Settings{}

	settings.LoadSettings() // #nosec G104

	return settings
}

func (s *Settings) SaveSettings() error {
	mutx.Lock()
	defer mutx.Unlock()

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	//save the file with read/write perms for admin only
	if err := util.WriteFile(platform.GetSettingsFile(), data, 0600); err != nil {
		return err
	}

	return nil
}

func (s *Settings) LoadSettings() error {
	mutx.RLock()
	defer mutx.RUnlock()

	data, err := os.ReadFile(platform.GetSettingsFile())
	if err != nil {
		return fmt.Errorf("unable to read settings: %w", err)
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		return err
	}

	return nil
}

func GetSessionCookie() string {
	sett := NewSettings()
	return "RADICAL_SESSION_ID=" + sett.Session.Secret + ";"
}
