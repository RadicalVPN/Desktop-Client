//go:build darwin

package main

import (
	"os"
)

func IsAdmin() bool {
	return os.Geteuid() == 0
}
