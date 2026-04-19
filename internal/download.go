package internal

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Release struct {
	IsLatest bool   `json:"isLatest"`
	Name     string `json:"name"`
}

func EnsureStoreDir(repo, version string) (string, error) {
	// Resolve base dir
	base := os.Getenv("WARES_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".local", "share", "wares")
	}

	// Split repo into owner/name
	parts := strings.Split(repo, "/")

	owner := parts[0]
	name := parts[1]

	// Build full path
	dir := filepath.Join(base, owner, name, version)

	// Create dir
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	return dir, nil
}

func BaseStoreDir() (string, error) {
	// Resolve base dir
	base := os.Getenv("WARES_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".local", "share", "wares")
	}

	return base, nil
}

func GetReleases(repo string) ([]Release, error) {
	out, err := exec.Command("gh", "release", "list", "--repo", repo, "--json", "name,isLatest").Output()
	if err != nil {
		return nil, err
	}

	var data []Release
	err = json.Unmarshal(out, &data)
	if err != nil {
		return nil, err
	}
	return data, nil

}

func GetLatest(repo string) (string, error) {
	data, err := GetReleases(repo)
	if err != nil {
		return "", err
	}

	for _, r := range data {
		if r.IsLatest && r.Name != "" {
			return r.Name, nil
		}
	}

	return "", nil

}

func Download(repo, release, pattern string) error {
	dir, err := EnsureStoreDir(repo, release)
	if err != nil {
		return err
	}
	command := exec.Command("gh", "release", "download", "--repo", repo, "--pattern", pattern, "--dir", dir, release)
	command.Stdout = os.Stdout
	err = command.Run()
	if err != nil {
		return err
	}

	return nil
}
