package cmd

import (
	"fmt"

	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

var (
	queryWare      bool
	queryBlueprint bool
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query information about an installed ware or blueprint",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if queryWare && queryBlueprint {
			fmt.Printf("%s Only one of --ware or --blueprint may be specified\n", internal.ErrText)
		}

		lock, err := internal.LoadLock()
		if err != nil {
			fmt.Printf("%s Error loading lockfile: %s\n", internal.ErrText, err)
		}

		if queryBlueprint {
			err := internal.QueryBlueprint(lock, name)
			if err != nil {
				fmt.Printf("%s Error querying blueprint: %s\n", internal.ErrText, err)
			}
		} else {
			err := internal.QueryWare(lock, name)
			if err != nil {
				fmt.Printf("%s Error querying ware: %s\n", internal.ErrText, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().BoolVar(&queryWare, "ware", false, "Query a ware")
	queryCmd.Flags().BoolVar(&queryBlueprint, "blueprint", false, "Query a blueprint")
}
