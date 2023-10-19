package service

import (
	"radicalvpnd/logger"
	"radicalvpnd/wireguard"
)

var log *logger.Logger

type Service struct {
	connected bool
}

func NewService() *Service {
	return &Service{}
}

func init() {
	log = logger.NewLogger("service")
}

func (s *Service) Connect(config string) {
	s.connected = true
	wireguardService := wireguard.NewWireguard()

	wireguardService.Connect("test configuration")
}

func (s *Service) Connected() bool {
	return s.connected
}
