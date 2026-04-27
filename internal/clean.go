package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CleanDir(repo, version string) error {
	dir, err := EnsureStoreDir(repo, version)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
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

func cleanWareVersions(repo, keepVersion string) error {
	base, err := BaseStoreDir()
	if err != nil {
		return err
	}

	parts := strings.Split(repo, "/")
	repoDir := filepath.Join(base, parts[0], parts[1])

	entries, err := os.ReadDir(repoDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		version := e.Name()

		if version == keepVersion {
			continue
		}

		fmt.Printf("%s Cleaning %s %s\n", CleanText, parts[1], version)
		err := CleanDir(repo, version)
		if err != nil {
			return err
		}
	}

	return nil
}

func Clean() error {
	lock, err := LoadLock()
	if err != nil {
		return err
	}

	for _, l := range lock.Wares {
		if err := cleanWareVersions(l.Repo, l.Version); err != nil {
			return err
		}
	}

	return nil
}
