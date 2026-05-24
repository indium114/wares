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

func ShellUpdate(dir string) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	cfg, err := LoadShellConfig(absDir)
	if err != nil {
		return err
	}

	lock, err := LoadShellLock(absDir)
	if err != nil {
		return err
	}

	changed := false

	for name, w := range cfg.Wares {
		fmt.Printf("%s %s %s ->", UpdateText, name, lock.Wares[name].Version)

		l, ok := lock.Wares[name]
		if !ok {
			l = LockedWare{}
		}

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

		if l.Version != latest {
			l = LockedWare{
				Repo:    w.Repo,
				Version: latest,
				Asset:   "",
				Digest:  "",
			}
			lock.Wares[name] = l
			changed = true
		}

		fmt.Printf("%s\n", l.Version)
	}

	for name, bp := range cfg.Blueprints {
		fmt.Printf("%s %s", UpdateText, name)

		repoDir, err := ensureBlueprintRepo(bp.Repo)
		if err != nil {
			return err
		}

		latest, err := resolveLatestCommit(repoDir)
		if err != nil {
			return err
		}

		locked := lock.Blueprints[name]
		if locked.Commit != latest {
			lock.Blueprints[name] = LockedBlueprint{
				Repo:        bp.Repo,
				Commit:      latest,
				BuiltCommit: locked.BuiltCommit,
			}
			changed = true
		}
	}

	if changed {
		return SaveShellLock(absDir, lock)
	}

	return nil
}
