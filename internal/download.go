package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Release struct {
	IsLatest bool   `json:"isLatest"`
	Name     string `json:"tagName"`
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
	out, err := exec.Command("gh", "release", "list", "--repo", repo, "--json", "tagName,isLatest").Output()
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

func GiteaGetLatest(host, repo string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/releases/latest", host, repo)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	type GiteaRelease struct {
		TagName string `json:"tag_name"`
	}

	var data GiteaRelease
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.TagName, nil
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

func downloadFile(downloadURL, dir, filename string) error {
	// Download file
	response, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// create destination file
	path := filepath.Join(dir, filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// copy response body to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func Download(repo, release, pattern, host string) error {
	dir, err := EnsureStoreDir(repo, release)
	if err != nil {
		return err
	}

	if host == "" || host == "https://github.com" {
		command := exec.Command("gh", "release", "download", "--repo", repo, "--pattern", pattern, "--dir", dir, release, "--clobber")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			return err
		}
	} else {
		type asset struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			UUID string `json:"uuid"`
			URL  string `json:"browser_download_url"`
		}
		type giteaRelease struct {
			Id     int     `json:"id"`
			Assets []asset `json:"assets"`
		}

		// get release id from tag
		url := fmt.Sprintf("%s/api/v1/repos/%s/releases/tags/%s", host, repo, release)
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var data giteaRelease
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}

		// get asset url by matching against pattern
		var id string
		for _, a := range data.Assets {
			if wildcardMatch(pattern, a.Name) {
				id = strconv.Itoa(a.Id)
			}
		}

		// download asset
		var downURL string
		var filename string
		for _, a := range data.Assets {
			if strconv.Itoa(a.Id) == id {
				downURL = a.URL
				filename = a.Name
			}
		}

		return downloadFile(downURL, dir, filename)
	}

	return nil
}
