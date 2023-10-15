//go:build darwin

package main

import (
	"os"
)

func isAdmin() bool {
	return os.Geteuid() == 0
}
