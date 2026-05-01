package internal

import (
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
	fullCmd := command + " " + pkg
	cmd := exec.Command(fullCmd)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func uninstallManagerOrphans() {}
