package cli

import (
	"bytes"
	"os/exec"
)

func Exec(binary string, args ...string) (string, error) {
	cmd := exec.Command(binary, args...)
	var outBuffer bytes.Buffer

	cmd.Stdout = &outBuffer
	cmd.Stderr = &outBuffer

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return outBuffer.String(), nil
}
