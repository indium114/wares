package cmd

import (
	"github.com/indium114/slag"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if queryWare && queryBlueprint {
			return slag.Err("Only one of --ware or --blueprint may be specified\n")
		}

		lock, err := internal.LoadLock()
		if err != nil {
			return slag.Err("Error loading lockfile: %s\n", err)
		}

		if queryBlueprint {
			err := internal.QueryBlueprint(lock, name)
			if err != nil {
				return slag.Err("Error querying blueprint: %s\n", err)
			}
		} else {
			err := internal.QueryWare(lock, name)
			if err != nil {
				return slag.Err("Error querying ware: %s\n", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().BoolVar(&queryWare, "ware", false, "Query a ware")
	queryCmd.Flags().BoolVar(&queryBlueprint, "blueprint", false, "Query a blueprint")
}
