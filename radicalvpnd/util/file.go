package util

import (
	"os"
	"runtime"

	"github.com/hectane/go-acl"
)

func FileChmod(file string, fileMode os.FileMode) error {
	if runtime.GOOS == "windows" {
		if err := acl.Chmod(file, fileMode); err != nil {
			os.Remove(file) // #nosec G104
			return err
		}
	} else {
		if err := os.Chmod(file, fileMode); err != nil {
			os.Remove(file) // #nosec G104
			return err
		}
	}
	return nil
}

func WriteFile(path string, data []byte, fileMode os.FileMode) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer file.Close()

	// Ensure file has correct permissions
	if err := FileChmod(path, fileMode); err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}
