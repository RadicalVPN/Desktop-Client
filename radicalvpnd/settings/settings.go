package settings

import (
	"encoding/json"
	"fmt"
	"radicalvpnd/platform"
	"radicalvpnd/util"
	"sync"
)

type Session struct {
	Session string
}

type Settings struct {
	Session Session
}

var mutx sync.RWMutex

func NewSettings() *Settings {
	return &Settings{}
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
