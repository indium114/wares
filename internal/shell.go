package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type ShellConfig struct {
	Wares      map[string]Ware      `yaml:"wares"`
	Blueprints map[string]Blueprint `yaml:"blueprint"`
}

func LoadShellConfig(dir string) (*ShellConfig, error) {
	path := filepath.Join(dir, "waresfile.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s waresfile.yaml not found in %s", ErrText, dir)
		}
		return nil, err
	}

	var cfg ShellConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Wares == nil {
		cfg.Wares = map[string]Ware{}
	}
	if cfg.Blueprints == nil {
		cfg.Blueprints = map[string]Blueprint{}
	}

	return &cfg, nil
}

func LoadShellLock(dir string) (*Lockfile, error) {
	path := filepath.Join(dir, "wares.lock")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Lockfile{
				Wares:      map[string]LockedWare{},
				Blueprints: map[string]LockedBlueprint{},
			}, nil
		}
		return nil, err
	}

	var lock Lockfile
	if err := yaml.Unmarshal(data, &lock); err != nil {
		return nil, err
	}

	if lock.Wares == nil {
		lock.Wares = map[string]LockedWare{}
	}
	if lock.Blueprints == nil {
		lock.Blueprints = map[string]LockedBlueprint{}
	}
	if lock.Managers == nil {
		lock.Managers = map[string][]string{}
	}

	return &lock, nil
}
