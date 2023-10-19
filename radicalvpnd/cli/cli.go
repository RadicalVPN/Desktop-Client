package cli

import "os/exec"

func Exec(binary string, args ...string) error {
	cmd := exec.Command(binary, args...)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
