package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func findManagerOrphans(cfg *Config, lock *Lockfile) map[string][]string {
	orphans := map[string][]string{}

	for managerName, lockedPkgs := range lock.Managers {
		configuredPkgs, configured := cfg.Managers[managerName]
		if !configured {
			// if entire manager was removed from config
			orphans[managerName] = lockedPkgs
			continue
		}

		configSet := make(map[string]bool, len(configuredPkgs))
		for _, pkg := range configuredPkgs {
			configSet[pkg] = true
		}

		for _, pkg := range lockedPkgs {
			if !configSet[pkg] {
				orphans[managerName] = append(orphans[managerName], pkg)
			}
		}
	}

	return orphans
}

func runManagerCommand(command, pkg string) error {
	cmd := exec.Command("sh", "-c", command, pkg)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func uninstallManagerOrphans(cfg *Config, lock *Lockfile) error {
	orphans := findManagerOrphans(cfg, lock)

	for managerName, pkgs := range orphans {
		// check that manager is configured
		settings, exists := cfg.Settings.Managers[managerName]
		if !exists {
			fmt.Printf("%s No settings for manager %s, skipping removal\n", WarnText, managerName)
			continue
		}

		// remove the pkg
		for _, pkg := range pkgs {
			fmt.Printf("%s Removing %s/%s\n", SyncText, managerName, pkg)
			if err := runManagerCommand(settings.Remove, pkg); err != nil {
				return err
			}

			// unlock pkg
			lockPkgs := lock.Managers[managerName]
			filtered := make([]string, 0, len(lockPkgs))
			for _, p := range lockPkgs {
				if p != pkg {
					filtered = append(filtered, p)
				}
			}
			lock.Managers[managerName] = filtered

			if len(filtered) == 0 {
				delete(lock.Managers, managerName)
			}
		}
	}

	return nil
}

func SyncManagers(cfg *Config, lock *Lockfile) (bool, error) {
	changed := false

	for managerName, pkgs := range cfg.Managers {
		settings, exists := cfg.Settings.Managers[managerName]
		if !exists {
			fmt.Printf("%s No settings for manager %s, skipping\n", WarnText, managerName)
			continue
		}

		if lock.Managers == nil {
			lock.Managers = map[string][]string{}
		}

		lockedSet := make(map[string]bool)
		for _, p := range lock.Managers[managerName] {
			lockedSet[p] = true
		}

		for _, pkg := range pkgs {
			if lockedSet[pkg] {
				continue
			}

			fmt.Printf("%s Installing %s/%s\n", SyncText, managerName, pkg)
			if err := runManagerCommand(settings.Install, pkg); err != nil {
				return false, err
			}

			lock.Managers[managerName] = append(lock.Managers[managerName], pkg)
			changed = true
		}
	}

	// remove orphans
	if err := uninstallManagerOrphans(cfg, lock); err != nil {
		return false, err
	}

	// remove manager from lockfile if no locked packages
	for managerName, pkgs := range lock.Managers {
		if len(pkgs) == 0 {
			delete(lock.Managers, managerName)
		}
	}

	return changed, nil
}
