package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/goccy/go-yaml"
	"io"
	"os"
	"path/filepath"
)

// MARK: Types
type Config struct {
	Wares map[string]Ware `yaml:"wares"`
}

type Ware struct {
	Repo           string `yaml:"repo"`
	Asset          string `yaml:"asset"`
	Name           string `yaml:"name"`
	Multiple       bool   `yaml:"multiple"`
	RemoveTopLevel bool   `yaml:"removetoplevel"`
}

type Lockfile struct {
	Wares map[string]LockedWare `yaml:"wares"`
}

type LockedWare struct {
	Repo    string `yaml:"repo"`
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

	return filepath.Join(dir, "config.yaml"), nil
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

func SaveLock(lock *Lockfile) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(dir, "pallet.lock")
	tmp := path + ".tmp"

	// ensure maps are not nil (prevents ugly YAML output)
	if lock.Wares == nil {
		lock.Wares = map[string]LockedWare{}
	}

	data, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}

	// write to temp file first
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}

	// atomic replace
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

func LinkWare(name, repo, version, asset string) error {
	// resolve store directory
	dir, err := EnsureStoreDir(repo, version)
	if err != nil {
		return err
	}

	target := filepath.Join(dir, asset)

	// validate target exists
	if _, err := os.Stat(target); err != nil {
		return err
	}

	// resolve symlink location
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	linkPath := filepath.Join(home, "Wares", name)

	// ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		return err
	}

	// remove existing symlink/file
	if err := os.Remove(linkPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// create new symlink
	if err := os.Symlink(target, linkPath); err != nil {
		return err
	}

	return nil
}

func ComputeDigest(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
