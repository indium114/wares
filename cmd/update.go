package cmd

import (
	"github.com/indium114/slag"
	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update packages in pallet.lock",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := internal.Update()
		if err != nil {
			return slag.Err("Failed to update: %s", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
