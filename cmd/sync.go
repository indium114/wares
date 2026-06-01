package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/indium114/slag"
	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

var clean bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync installed packages with config",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if config directory is a git repo, and if it has unstaged changes
		configDir, _ := internal.ConfigDir()
		command := exec.Command("git", "-C", configDir, "status", "--porcelain")
		out, err := command.Output()
		if string(out) != "" && string(out) != "fatal: not a git repository (or any of the parent directories): .git" {
			slag.Warn("Git tree dirty (remember to commit your changes)\n")
		}

		err = internal.Sync(clean)
		if err != nil {
			fmt.Print(slag.Err("Failed to sync: %s", err).Error())
		}

		slag.Log("Marking all files in ~/Wares as executable\n")
		err = filepath.Walk(os.ExpandEnv("$HOME/Wares"), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				mode := info.Mode()
				return os.Chmod(path, mode|0111)
			}

			return nil
		})

		slag.Log("Marking all files in /Wares as executable\n")
		err = filepath.Walk("/Wares", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				mode := info.Mode()
				return os.Chmod(path, mode|0111)
			}

			return nil
		})

		if err != nil {
			return slag.Err("Failed to mark files as executable: %s", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVar(&clean, "clean", false, "Rebuild all blueprints from scratch")
}
