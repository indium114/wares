package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

type ShellConfig struct {
	Wares      map[string]Ware      `yaml:"wares"`
	Blueprints map[string]Blueprint `yaml:"blueprints"`
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
		fmt.Printf("%s %s %s -> ", UpdateText, name, lock.Wares[name].Version)

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
		fmt.Printf("%s %s\n", UpdateText, name)

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
				Artifacts:   locked.Artifacts,
			}
			changed = true
		}
	}

	if changed {
		return SaveShellLock(absDir, lock)
	}

	return nil
}

func ShellSync(dir string) error {
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

	shellDir := shellWaresDir(absDir)
	changed := false

	var results []syncResult

	for name, w := range cfg.Wares {
		fmt.Printf("%s Installing %s\n", SyncText, name)

		l, ok := lock.Wares[name]
		if !ok || l.Version == "" {
			return fmt.Errorf("%s %s not locked yet, run 'wares shell --update' first", ErrText, name)
		}

		if err := CleanDir(w.Repo, l.Version); err != nil {
			return err
		}

		if err := Download(w.Repo, l.Version, w.Asset, w.Host); err != nil {
			return err
		}

		path, err := ResolveDownloadedPath(w.Repo, l.Version, w.Asset, w.Multiple)
		if err != nil {
			return err
		}

		digest, err := ComputeDigest(path)
		if err != nil {
			return err
		}

		if l.Digest == "" {
			l.Digest = "sha256:" + digest
			l.Asset = filepath.Base(path)
			lock.Wares[name] = l
			changed = true
		} else {
			expected := strings.TrimPrefix(l.Digest, "sha256:")
			if expected != digest {
				return fmt.Errorf("%s digest mismatch", name)
			}
		}

		if archive, kind := IsArchive(path); archive {
			if err := Extract(path, filepath.Dir(path), kind, w.RemoveTopLevel); err != nil {
				return err
			}
		}

		var linkSource string
		if archive, _ := IsArchive(path); archive {
			linkSource = w.Name
		} else {
			linkSource = l.Asset
		}

		if w.Multiple {
			linkSource = filepath.Dir(linkSource)
		}

		results = append(results, syncResult{
			name:       name,
			repo:       w.Repo,
			version:    l.Version,
			linkSource: linkSource,
		})
	}

	for _, r := range results {
		linkPath := filepath.Join(shellDir, r.name)
		os.Remove(linkPath)

		if err := shellSymlinkWare(r.name, r.repo, r.version, r.linkSource, shellDir); err != nil {
			return err
		}
	}

	for name, bp := range cfg.Blueprints {
		fmt.Printf("%s Building %s\n", SyncText, name)

		repoDir, err := ensureBlueprintRepo(bp.Repo)
		if err != nil {
			return err
		}

		locked, ok := lock.Blueprints[name]
		if !ok || locked.Commit == "" {
			return fmt.Errorf("%s %s not locked yet, run 'wares shell --update' first", ErrText, name)
		}

		needRebuild := locked.BuiltCommit != locked.Commit || locked.Repo != bp.Repo
		if needRebuild {
			if err := buildBlueprint(repoDir, locked.Commit, bp.Steps); err != nil {
				return err
			}

			lock.Blueprints[name] = LockedBlueprint{
				Repo:        bp.Repo,
				Commit:      locked.Commit,
				BuiltCommit: locked.Commit,
				Artifacts:   bp.Artifacts,
			}
			changed = true
		}

		for _, artifact := range bp.Artifacts {
			if err := shellSymlinkBlueprint(artifact, repoDir, shellDir); err != nil {
				return err
			}
		}
	}

	fmt.Printf("%s Marking all files in %s as executable\n", LogText, shellDir)
	filepath.Walk(shellDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		return os.Chmod(path, info.Mode()|0111)
	})

	if changed {
		return SaveShellLock(absDir, lock)
	}

	return nil
}
