package cmd

import (
	"fmt"

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
			fmt.Println("%s Failed to update: %s", internal.ErrText, err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
