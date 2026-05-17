package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func removeLink(name string, system bool) error {
	waresDir, err := WaresDir(system)
	if err != nil {
		return err
	}

	linkPath := filepath.Join(waresDir, name)

	err = os.Remove(linkPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func Update() error {
	cfg, err := LoadConfig()
	lock, err := LoadLock()
	if err != nil {
		return err
	}

	for name, w := range cfg.Wares {
		fmt.Printf("%s %s %s -> ", UpdateText, name, lock.Wares[name].Version)

		var latest string
		if w.Host == "" || w.Host == "https://github.com" {
			latest, err = GetLatest(w.Repo)
			if err != nil {
				return err
			}
		} else {
			latest, err = GiteaGetLatest(w.Host, w.Repo)
			if err != nil {
				return err
			}

		}

		l, ok := lock.Wares[name]
		if !ok {
			l = LockedWare{}
		}

		// update if changes
		if l.Version != latest {
			l.Version = latest
			l.Digest = ""
			lock.Wares[name] = l
		}

		fmt.Printf("%s\n", l.Version)
	}

	for name, bp := range cfg.Blueprints {
		fmt.Printf("%s %s\n", UpdateText, name)

		// pull latest
		repoDir, err := ensureBlueprintRepo(bp.Repo)
		if err != nil {
			return err
		}

		// get commit
		latest, err := resolveLatestCommit(repoDir)
		if err != nil {
			return err
		}

		// lock
		locked := lock.Blueprints[name]
		if locked.Commit != latest {
			lock.Blueprints[name] = LockedBlueprint{
				Repo:   bp.Repo,
				Commit: latest,
				System: bp.System,
			}
		}
	}

	if err := UpdateManagers(cfg, lock); err != nil {
		return err
	}

	return SaveLock(lock)
}
