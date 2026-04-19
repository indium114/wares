package internal

import (
	"os"
	"path/filepath"
	"strings"
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

func RemoveStore(repo string) error {
	base, err := BaseStoreDir()
	if err != nil {
		return err
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return err
	}

	path := filepath.Join(base, parts[0], parts[1])

	return os.RemoveAll(path)
}
