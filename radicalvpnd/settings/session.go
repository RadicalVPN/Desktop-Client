package settings

type Session struct {
	Secret string `json:"secret"`
}

func (s *Settings) SetSession(secret string) {
	s.LoadSettings()

	s.Session.Secret = secret

	s.SaveSettings()
}
