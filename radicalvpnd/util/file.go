package util

import (
	"os"
	"runtime"

	"github.com/hectane/go-acl"
)

func FileChmod(file string, fileMode os.FileMode) error {
	if runtime.GOOS == "windows" {
		if err := acl.Chmod(file, fileMode); err != nil {
			os.Remove(file)
			return err
		}
	} else {
		if err := os.Chmod(file, fileMode); err != nil {
			os.Remove(file)
			return err
		}
	}
	return nil
}

func WriteFile(file string, data []byte, fileMode os.FileMode) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	// Ensure file has correct permissions
	if err := FileChmod(file, fileMode); err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}
