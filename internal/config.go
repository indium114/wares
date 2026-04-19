package internal

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
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

// MARK: Parsing

func readYamlFile(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}

func LoadConfig() (*Config, error) {
	var cfg Config

	path, _ := WaresFile()

	err := readYamlFile(path, &cfg)
	if err != nil {
		if os.IsNotExist(err) {
			// return empty config
			return &Config{
				Wares: map[string]Ware{},
			}, nil
		}
		return nil, err
	}

	if cfg.Wares == nil {
		cfg.Wares = map[string]Ware{}
	}

	return &cfg, nil
}

func LoadLock() (*Lockfile, error) {
	var lock Lockfile

	path, _ := LockFile()

	err := readYamlFile(path, &lock)
	if err != nil {
		if os.IsNotExist(err) {
			// return empty config
			return &Lockfile{
				Wares: map[string]LockedWare{},
			}, nil
		}
		return nil, err
	}

	if lock.Wares == nil {
		lock.Wares = map[string]LockedWare{}
	}

	return &lock, nil
}
