package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolveDownloadedPath(repo, version, pattern string, multiple bool) (string, error) {
	dir, err := EnsureStoreDir(repo, version)
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var files []os.DirEntry

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e)
	}

	if len(files) == 0 {
		fmt.Printf("%s No downloaded artifact in %s\n", ErrText, dir)
		return "", err
	}

	if len(files) > 1 {
		names := make([]string, 0, len(files))
		for _, f := range files {
			names = append(names, f.Name())
		}
		if multiple {
			return filepath.Join(dir, files[0].Name()), nil
		} else {
			return "", fmt.Errorf("%s Multiple artifacts found in %s: %v", ErrText, dir, names)
		}
	}

	return filepath.Join(dir, files[0].Name()), nil
}

func findOrphans() []string {
	cfg, err := LoadConfig()
	lock, err := LoadLock()
	if err != nil {
		return []string{}
	}

	var orphans []string

	for name := range lock.Wares {
		if _, ok := cfg.Wares[name]; !ok {
			orphans = append(orphans, name)
		}
	}

	return orphans
}

func UninstallOrphans() error {
	lock, err := LoadLock()
	if err != nil {
		return err
	}

	orphans := findOrphans()

	for _, name := range orphans {
		fmt.Printf("%s Removing %s\n", SyncText, name)

		l := lock.Wares[name]

		// remove symlink
		if err := removeLink(name); err != nil {
			return err
		}

		// remove stored files
		if err := RemoveStore(l.Repo); err != nil {
			return err
		}

		//  remove from lockfile
		delete(lock.Wares, name)

		SaveLock(lock)
	}

	return nil
}

func Sync() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	lock, err := LoadLock()
	if err != nil {
		return err
	}

	changed := false

	for name, w := range cfg.Wares {
		// Remove old version
		if err := removeLink(name); err != nil {
			return err
		}

		fmt.Printf("%s Installing %s\n", SyncText, name)

		l, ok := lock.Wares[name]

		// Missing lock entry → resolve latest
		if !ok || l.Version == "" {
			rel, err := GetLatest(w.Repo)
			if err != nil {
				return fmt.Errorf("%s %s: latest release: %w", ErrText, name, err)
			}

			l = LockedWare{
				Repo:    w.Repo,
				Version: rel,
				Asset:   "", // will be filled after download
				Digest:  "",
			}

			lock.Wares[name] = l
			changed = true
		}

		if err := CleanDir(w.Repo, l.Version); err != nil {
			return err
		}

		// Download
		if err := Download(w.Repo, l.Version, w.Asset); err != nil {
			return err
		}

		// Resolve installed file path
		path, err := ResolveDownloadedPath(w.Repo, l.Version, w.Asset, w.Multiple)
		if err != nil {
			return fmt.Errorf("%s: resolve path: %w", name, err)
		}

		// Compute digest
		digest, err := ComputeDigest(path)
		if err != nil {
			return fmt.Errorf("%s: digest: %w", name, err)
		}

		// Lock or verify
		if l.Digest == "" {
			l.Digest = "sha256:" + digest
			l.Asset = filepath.Base(path)
			lock.Wares[name] = l
			changed = true
		} else {
			expected := strings.TrimPrefix(l.Digest, "sha256:")
			if expected != digest {
				return fmt.Errorf("%s: digest mismatch", name)
			}
		}

		// Extract if .tar.gz
		if archive, kind := IsArchive(path); archive == true {
			if err := Extract(path, filepath.Dir(path), kind, (w.RemoveTopLevel || false)); err != nil {
				return err
			}
		}

		// Symlink
		var linkSource string
		if archive, _ := IsArchive(path); archive == true {
			linkSource = w.Name
		} else {
			linkSource = l.Asset
		}

		if w.Multiple {
			linkSource = filepath.Dir(linkSource)
		}

		if err := LinkWare(name, w.Repo, l.Version, linkSource); err != nil {
			return fmt.Errorf("%s: link: %w", name, err)
		}
	}

	// Sync native managers
	if changed, err = SyncManagers(cfg, lock); err != nil {
		return err
	}

	// Remove orphaned packages
	UninstallOrphans()

	// Persist lockfile if modified
	if changed {
		if err := SaveLock(lock); err != nil {
			return err
		}
	}

	return nil
}
