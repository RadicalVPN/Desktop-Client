package service

import (
	"radicalvpnd/logger"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("service")
}
