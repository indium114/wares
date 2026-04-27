package cmd

import (
	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove old versions of packages",
	Run: func(cmd *cobra.Command, args []string) {
		internal.UninstallOrphans()
		internal.Clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
