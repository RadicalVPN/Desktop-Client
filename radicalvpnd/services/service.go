package service

import (
	"radicalvpnd/logger"
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
