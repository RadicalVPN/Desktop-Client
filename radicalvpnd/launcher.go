package main

import (
	"fmt"
	"os"
	"radicalvpnd/logger"
	"radicalvpnd/platform"
	"radicalvpnd/protocol"
	"radicalvpnd/util"
	"radicalvpnd/version"
	"runtime"
	"strconv"
	"strings"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("launch")
}

func Launch() {
	logger.Init(platform.GetLogFilePath())

	log.Info("Starting RadicalVPN Daemon..", fmt.Sprintf(" [%s,%s]", runtime.GOOS, runtime.GOARCH))
	log.Info("Version: ", version.GetVersion())
	log.Info("Commit Hash: ", version.GetCommitHash())

	if version.IsNightly() && version.IsNightlyOutdated() {
		log.Warning("Nightly Build of Daemon is outdated!")
	}

	if version.IsRelease() && version.IsReleaseOutdated() {
		log.Warning("Production Build of Daemon is outdated!")
	}

	log.Info(fmt.Sprintf("Args: %s", os.Args))
	log.Info(fmt.Sprintf("PID : %d PPID: %d", os.Getpid(), os.Getppid()))
	log.Info(fmt.Sprintf("Arch: %d bit", strconv.IntSize))

	defer func() {
		log.Info("RadicalVPN Daemon stopped.")
	}()

	log.Info("Initializing platform specific variables..")
	platform.Init()

	//check if the daemon is started as admin/root
	if !IsAdmin() {
		log.Warning(strings.Repeat("-", 48))
		log.Warning("! NO ADMIN USER !")
		log.Warning("RadicalVPN Daemon must be started as admin/root!")
		log.Warning(strings.Repeat("-", 48))
		os.Exit(1)
	}

	var secret string
	envSecret, envSecretPresent := os.LookupEnv("RADICALVPND_SECRET")
	if envSecretPresent {
		secret = envSecret
	} else {
		generatedSecret, err := util.RandomString(32)
		if err != nil {
			log.Error("Failed to generate random string: ", err)
			os.Exit(1)
		}

		secret = generatedSecret
	}

	portChannel := make(chan string, 1)
	go func() {
		//this will block until we got a port back from the protocol
		port := <-portChannel

		log.Info("Found listening port in port channel: ", port)

		//write port to service file
		file, err := os.Create(platform.GetServiceFIle())
		if err != nil {
			log.Error("Failed to open service file: ", err)
			os.Exit(1)
		}

		defer file.Close()

		log.Info("Writing port and secret to service file..")
		file.WriteString(fmt.Sprintf("%s|%s", port, secret))
	}()

	if err := PrepareRun(); err != nil {
		log.Error("Failed to prepare run: ", err)
		os.Exit(1)
	}

	protocol := protocol.NewProtocol(secret)
	protocol.Start(portChannel)
}
