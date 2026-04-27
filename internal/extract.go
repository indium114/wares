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

func Extract(archive, dir string, removeTopLevel bool) error {
	var removeTopLevelArg string
	if removeTopLevel {
		removeTopLevelArg = "--strip-components=1"
	} else {
		removeTopLevelArg = ""
	}

	command := exec.Command("tar", "xvf", archive, "--directory", dir, removeTopLevelArg)
	err := command.Run()
	if err != nil {
		return err
	}

	return nil
}
