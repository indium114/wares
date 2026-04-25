package internal

import (
	"os/exec"
	"strings"
)

func IsArchive(name string) bool {
	if strings.HasSuffix(name, ".gz") {
		return true
	}

	return false
}

func Extract(archive, dir string) error {
	command := exec.Command("tar", "xvf", archive, "--directory", dir)
	err := command.Run()
	if err != nil {
		return err
	}

	return nil
}
