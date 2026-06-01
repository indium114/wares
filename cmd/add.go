package cmd

import (
	"fmt"

	"github.com/indium114/slag"
	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a ware configuration from warehouse to config.yaml",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgName := args[0]

		slag.Add("Adding ware %s from warehouse to config\n", pkgName)

		cfg, err := internal.LoadConfig()
		if err != nil {
			return err
		}

		// check that platform and arch settings are set
		if cfg.Settings.Platform == "" {
			return slag.Err("settings:platform not set in config.yaml (must be 'linux' or 'darwin')\n")
		}

		if cfg.Settings.Arch == "" {
			return slag.Err("settings:arch not set in config.yaml (must be 'x86_64' or 'aarch64')\n")
		}

		// validate platform and arch settings
		validPlatforms := map[string]bool{"linux": true, "darwin": true}
		if !validPlatforms[cfg.Settings.Platform] {
			return slag.Err("Invalid platform %q (must be 'linux' or 'darwin')\n", cfg.Settings.Platform)
		}

		validArches := map[string]bool{"x86_64": true, "aarch64": true}
		if !validArches[cfg.Settings.Arch] {
			return slag.Err("Invalid platform %q (must be 'x86_64' or 'aarch64')\n", cfg.Settings.Arch)
		}

		// skip if user has already configured package
		if _, exists := cfg.Wares[pkgName]; exists {
			slag.Warn("%s ware %q is already in config, skipping\n", pkgName)
			return nil
		}

		// get the package config from the warehouse
		ware, err := internal.FetchFromWarehouse(cfg.Settings.Platform, cfg.Settings.Arch, pkgName)
		if err != nil {
			if err.Error() == fmt.Sprintf("%s ware %q not found in warehouse", internal.ErrText, pkgName) {
				slag.Err("ware %q not found in the warehouse\n", pkgName)
				slag.Hint("You can contribute a package configuration at https://github.com/wares-pkg/warehouse\n")
				return nil
			}
			return slag.Err("%w", err)
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
