package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func removeLink(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	linkPath := filepath.Join(home, "Wares", name)

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
		fmt.Printf("%s %s\n", UpdateText, name)

		latest, err := GetLatest(w.Repo)
		if err != nil {
			return err
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
	}

	return SaveLock(lock)
}
