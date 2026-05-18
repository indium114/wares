package internal

import (
	"os"
	"os/exec"
	"strings"
)

func Sudo(args ...string) error {
	cmd := exec.Command("sh", "-c", "sudo "+strings.Join(args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
