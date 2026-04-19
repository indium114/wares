package internal

import (
	"os"
	"path/filepath"
)

// MARK: Types
type Config struct {
	Wares map[string]Ware `yaml:"wares"`
}

type Ware struct {
	Repo  string `yaml:"repo"`
	Asset string `yaml:"asset"`
	Name  string `yaml:"name"`
}

type Lockfile struct {
	Wares map[string]LockedWare `yaml:"wares"`
}

type LockedWare struct {
	Version string `yaml:"version"`
	Asset   string `yaml:"asset"`
	Digest  string `yaml:"digest"`
}

// MARK: Path functions
func ConfigDir() (string, error) {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "wares"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "wares"), nil
}

func WaresFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "wares.yaml"), nil
}

func LockFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "pallet.lock"), nil
}
