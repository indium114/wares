package internal

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-yaml"
)

func FetchFromWarehouse(platform, arch, pkgName string) (*Ware, error) {
	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/wares-pkg/warehouse/refs/heads/main/%s/%s/%s.yaml",
		platform, arch, pkgName,
	)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s ware %q not found in warehouse", ErrText, pkgName)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s Failed to fetch: HTTP %d", ErrText, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s Failed to read response: %w", ErrText, err)
	}

	var ware Ware
	if err := yaml.Unmarshal(data, &ware); err != nil {
		return nil, fmt.Errorf("%s Failed to parse package config: %w", ErrText, err)
	}

	return &ware, nil
}
