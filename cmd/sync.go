package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync installed packages with config",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := internal.Sync()
		if err != nil {
			fmt.Printf("%s Failed to sync: %s", internal.ErrText, err)
			return err
		}

		fmt.Printf("%s Marking all files in ~/Wares as executable\n", internal.LogText)
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

		if err != nil {
			fmt.Printf("%s Failed to mark files as executable: %s", internal.ErrText, err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
