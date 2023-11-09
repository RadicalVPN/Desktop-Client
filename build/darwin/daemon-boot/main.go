package main

import (
	"C"
	"os/exec"
)
import "log"

const daemon = "/Applications/RadicalVPN.app/Contents/MacOS/RadicalVPN Daemon"

func main() {

	println("Launching daemon...")
	cmd := exec.Command(daemon)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}