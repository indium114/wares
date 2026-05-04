package cmd

import (
	"fmt"

	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a ware configuration from warehouse to config.yaml",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgName := args[0]

		fmt.Printf("%s Adding ware %s from warehouse to config", internal.AddText, pkgName)

		cfg, err := internal.LoadConfig()
		if err != nil {
			return err
		}

		// check that platform and arch settings are set
		if cfg.Settings.Platform == "" {
			fmt.Printf("%s settings:platform not set in config.yaml (must be 'linux' or 'darwin')\n", internal.ErrText)
			return err
		}

		if cfg.Settings.Arch == "" {
			fmt.Printf("%s settings:arch not set in config.yaml (must be 'x86_64' or 'aarch64')\n", internal.ErrText)
			return err
		}

		// validate platform and arch settings
		validPlatforms := map[string]bool{"linux": true, "darwin": true}
		if !validPlatforms[cfg.Settings.Platform] {
			fmt.Printf("%s Invalid platform %q (must be 'linux' or 'darwin')\n", internal.ErrText, cfg.Settings.Platform)
		}

		validArches := map[string]bool{"x86_64": true, "aarch64": true}
		if !validArches[cfg.Settings.Arch] {
			fmt.Printf("%s Invalid platform %q (must be 'x86_64' or 'aarch64')\n", internal.ErrText, cfg.Settings.Arch)
		}

		// skip if user has already configured package
		if _, exists := cfg.Wares[pkgName]; exists {
			fmt.Printf("%s ware %q is already in config, skipping\n", internal.WarnText, pkgName)
			return nil
		}

		// get the package config from the warehouse
		ware, err := internal.FetchFromWarehouse(cfg.Settings.Platform, cfg.Settings.Arch, pkgName)
		if err != nil {
			if err.Error() == fmt.Sprintf("%s ware %q not found in warehouse", internal.ErrText, pkgName) {
				fmt.Printf("%s ware %q not found in the warehouse\n", internal.ErrText, pkgName)
				fmt.Printf("%s You can contribute a package configuration at https://github.com/wares-pkg/warehouse\n", internal.HintText)
				return nil
			}
			return fmt.Errorf("%s %w", internal.ErrText, err)
		}

		// add the ware to the config
		cfg.Wares[pkgName] = *ware

		// save config
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
