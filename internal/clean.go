package internal

import (
	"os"
	"path/filepath"
)

func CleanDir(repo, version string) error {
	dir, err := EnsureStoreDir(repo, version)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}
