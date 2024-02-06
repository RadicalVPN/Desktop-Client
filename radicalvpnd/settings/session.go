package settings

import "fmt"

type Session struct {
	Secret string `json:"secret"`
}

func (s *Settings) SetSession(secret string) {
	s.Session.Secret = secret

	err := s.SaveSettings()
	if err != nil {
		fmt.Printf("Failed to save session: %s\n", err)
	}
}
