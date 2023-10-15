package main

import (
	"fmt"
	"os"
	"radicalvpnd/logger"
	"radicalvpnd/protocol"
	"radicalvpnd/util"
	"runtime"
	"strconv"
	"strings"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("launch")
}

func Launch() {
	log.Info("Starting RadicalVPN Daemon..", fmt.Sprintf(" [%s,%s]", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Args: %s", os.Args))
	log.Info(fmt.Sprintf("PID : %d PPID: %d", os.Getpid(), os.Getppid()))
	log.Info(fmt.Sprintf("Arch: %d bit", strconv.IntSize))

	defer func() {
		log.Info("RadicalVPN Daemon stopped.")
	}()

	//check if the daemon is started as admin/root
	if isAdmin() == false {
		log.Warning(strings.Repeat("-", 48))
		log.Warning("! NO ADMIN USER !")
		log.Warning("RadicalVPN Daemon must be started as admin/root!")
		log.Warning(strings.Repeat("-", 48))
		os.Exit(1)
	}

	secret := util.RandomString(32)

	log.Debug("Secret: " + secret)

	protocol := protocol.NewProtocol(secret)
	protocol.Start()

}
