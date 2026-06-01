package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/indium114/slag"
	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

var (
	shellUpdate bool
	shellClean  bool
)

var shellCmd = &cobra.Command{
	Use:   "shell [dir]",
	Short: "Enter a waresfile.yaml shell",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."

		if len(args) > 0 {
			dir = args[0]
		}

		absDir, err := filepath.Abs(dir)
		if err != nil {
			return slag.Err("Failed to resolve absolute directory: %s", err)
		}

		if shellUpdate {
			if err := internal.ShellUpdate(absDir); err != nil {
				return slag.Err("Failed to update wares.lock: %s", err)
			}
		}

		if err := internal.ShellSync(absDir, shellClean); err != nil {
			return slag.Err("Failed to sync: %s", err)
		}

		shellDir := filepath.Join(absDir, ".wares")
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		slag.Hint("Entering wares shell\n")

		newEnv := os.Environ()
		for i, e := range newEnv {
			if strings.HasPrefix(e, "PATH=") {
				newEnv[i] = "PATH=" + shellDir + ":" + e[len("PATH"):]
				break
			}
		}
		newEnv = append(newEnv, "WARES_SHELL_ACTIVE=true")

		sh := exec.Command(shell)
		sh.Env = newEnv
		sh.Stdin = os.Stdin
		sh.Stdout = os.Stdout
		sh.Stderr = os.Stderr

		if err := sh.Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().BoolVar(&shellUpdate, "update", false, "Update wares.lock before entering the shell")
	shellCmd.Flags().BoolVar(&shellClean, "clean", false, "Rebuild all blueprints from scratch")
}
