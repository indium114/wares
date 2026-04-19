package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolveDownloadedPath(repo, version, pattern string) (string, error) {
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
		fmt.Println("%s No downloaded artifact in %s", ErrText, dir)
		return "", err
	}

	if len(files) > 1 {
		names := make([]string, 0, len(files))
		for _, f := range files {
			names = append(names, f.Name())
		}
		return "", fmt.Errorf("%s Multiple artifacts found in %s: %v", ErrText, dir, names)
	}

	return filepath.Join(dir, files[0].Name()), nil
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
		fmt.Printf("%s %s\n", SyncText, name)

		l, ok := lock.Wares[name]

		// Missing lock entry → resolve latest
		if !ok || l.Version == "" {
			rel, err := GetLatest(w.Repo)
			if err != nil {
				return fmt.Errorf("%s %s: latest release: %w", ErrText, name, err)
			}

			l = LockedWare{
				Version: rel,
				Asset:   "", // will be filled after download
				Digest:  "",
			}

			lock.Wares[name] = l
			changed = true
		}

		// Download
		if err := Download(w.Repo, l.Version, w.Asset); err != nil {
			fmt.Printf("%s Repo: %s, Version: %s, Asset: %s", ErrText, w.Repo, l.Version, w.Asset)
			return fmt.Errorf("%s: download: %w", name, err)
		}

		// Resolve installed file path
		path, err := ResolveDownloadedPath(w.Repo, l.Version, w.Asset)
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

		// Symlink
		if err := LinkWare(name, w.Repo, l.Version, l.Asset); err != nil {
			return fmt.Errorf("%s: link: %w", name, err)
		}
	}

	// Persist lockfile if modified
	if changed {
		if err := SaveLock(lock); err != nil {
			return err
		}
	}

	return nil
}
