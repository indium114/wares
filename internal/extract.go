package internal

import (
	"os/exec"
	"strings"
)

func IsArchive(name string) (bool, string) {
	if strings.HasSuffix(name, ".gz") {
		return true, "gz"
	} else if strings.HasSuffix(name, ".zip") {
		return true, "zip"
	}

	return false, ""
}

func Extract(archive, dir, kind string, removeTopLevel bool) error {
	var removeTopLevelArg string
	if removeTopLevel {
		if kind == "gz" {
			removeTopLevelArg = "--strip-components=1"
		} else {
			removeTopLevelArg = "-j"
		}
	} else {
		removeTopLevelArg = ""
	}

	var command *exec.Cmd
	if kind == "gz" {
		command = exec.Command("tar", "xvf", archive, "--directory", dir, removeTopLevelArg)
	} else {
		command = exec.Command("unzip", removeTopLevelArg, "-d", dir, archive)
	}
	err := command.Run()
	if err != nil {
		return err
	}

	return nil
}
