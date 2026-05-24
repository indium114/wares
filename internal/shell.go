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

func SaveShellLock(dir string, lock *Lockfile) error {
	path := filepath.Join(dir, "wares.lock")
	tmp := path + ".tmp"

	if lock.Wares == nil {
		lock.Wares = map[string]LockedWare{}
	}
	if lock.Blueprints == nil {
		lock.Blueprints = map[string]LockedBlueprint{}
	}
	if lock.Managers == nil {
		lock.Managers = map[string][]string{}
	}

	data, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}

	return nil
}

func shellWaresDir(absDir string) string {
	return filepath.Join(absDir, ".wares")
}

func shellSymlinkWare(name, repo, version, linkSource, shellDir string) error {
	storeDir, err := EnsureStoreDir(repo, version)
	if err != nil {
		return err
	}

	target := filepath.Join(storeDir, linkSource)

	if _, err := os.Stat(target); err != nil {
		return err
	}

	if err := os.MkdirAll(shellDir, 0o755); err != nil {
		return err
	}

	linkPath := filepath.Join(shellDir, name)

	os.Remove(linkPath)

	return os.Symlink(target, linkPath)
}

func shellSymlinkBlueprint(artifact, repoDir string, shellDir string) error {
	src := filepath.Join(repoDir, artifact)

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("%s artifact %s not found", ErrText, artifact)
	}

	if err := os.MkdirAll(shellDir, 0o755); err != nil {
		return err
	}

	linkPath := filepath.Join(shellDir, filepath.Base(artifact))

	os.Remove(linkPath)

	return os.Symlink(src, linkPath)
}
