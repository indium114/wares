package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/indium114/wares/internal"
	"github.com/spf13/cobra"
)

func checkGhCli() bool {
	command := exec.Command("gh", "auth", "status", "--json", "hosts")
	out, err := command.Output()
	if err != nil {
		return false
	}

	type GhAuthStatus struct {
		Hosts map[string][]struct {
			State  string `json:"state"`
			Active bool   `json:"active"`
			Host   string `json:"host"`
			Login  string `json:"login"`
		} `json:"hosts"`
	}

	var status GhAuthStatus
	if err := json.Unmarshal(out, &status); err != nil {
		return false
	}

	hosts, ok := status.Hosts["github.com"]
	if !ok {
		return false
	}

	for _, h := range hosts {
		if h.Active && h.State == "success" {
			return true
		}
	}

	return false
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check that all of wares' prerequisites are met",
	RunE: func(cmd *cobra.Command, args []string) error {
		// MARK: Check PATH for Wares directory
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		waresBin := filepath.Join(home, "Wares")
		if _, err := os.Stat(waresBin); os.IsNotExist(err) {
			fmt.Printf("%s %s does not exist, creating\n", internal.WarnText, waresBin)
		}
		err = os.MkdirAll(waresBin, 0o755)
		if err != nil {
			fmt.Printf("%s Failed to create wares bin dir: %s\n", internal.ErrText, err)
			return err
		}
		fmt.Printf("%s %s exists\n", internal.OkText, waresBin)

		path := os.Getenv("PATH")
		if !strings.Contains(path, waresBin) {
			fmt.Printf("%s %s not in PATH\n", internal.ErrText, waresBin)
			return err
		}
		fmt.Printf("%s %s is in PATH\n", internal.OkText, waresBin)

		// MARK: Check for Wares config directory
		waresConfig, _ := internal.ConfigDir()
		if _, err := os.Stat(waresConfig); os.IsNotExist(err) {
			fmt.Printf("%s %s does not exist, creating\n", internal.WarnText, waresConfig)
		}
		err = os.MkdirAll(waresConfig, 0o755)
		if err != nil {
			fmt.Printf("%s Failed to create wares config dir: %s\n", internal.ErrText, err)
		}
		fmt.Printf("%s %s exists\n", internal.OkText, waresConfig)

		// MARK: Check for authenticated gh CLI
		if status := checkGhCli(); status == true {
			fmt.Printf("%s Logged into GitHub CLI", internal.OkText)
		} else {
			fmt.Printf("%s Not logged into GitHub CLI", internal.ErrText)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
