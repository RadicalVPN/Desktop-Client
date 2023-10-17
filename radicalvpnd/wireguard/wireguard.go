package wireguard

import (
	"radicalvpnd/logger"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("wireguard")
}

type Wireguard struct {
	binary     string
	configFile string
}
