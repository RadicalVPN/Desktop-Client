package wireguard

import (
	"radicalvpnd/logger"
)

type Wireguard struct {
}

func NewWireguard() *Wireguard {
	return &Wireguard{}
}

var log *logger.Logger

func init() {
	log = logger.NewLogger("wireguard")

}

func (wg *Wireguard) Connect(config string) {
	log.Info("Connecting to wireguard..")
	wg.start(config)
}
