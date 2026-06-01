package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/indium114/slag"
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

func checkCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		return false
	}

	return true
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
			slag.Warn("%s does not exist, creating\n", waresBin)
		}
		err = os.MkdirAll(waresBin, 0o755)
		if err != nil {
			fmt.Print(slag.Err("Failed to create wares bin dir: %s\n", err).Error())
		}
		sysWaresBin := "/Wares"
		if _, err := os.Stat(sysWaresBin); os.IsNotExist(err) {
			slag.Warn("%s does not exist, creating\n", sysWaresBin)
		}
		err = internal.Sudo("mkdir", "-p", "/Wares")
		if err != nil {
			fmt.Print(slag.Err("Failed to create system wares bin dir: %s\n", err).Error())
		}
		err = internal.Sudo("chown", os.Getenv("USER"), "/Wares")
		if err != nil {
			fmt.Print(slag.Err("%s Failed to chown system wares bin dir: %s\n", internal.ErrText, err).Error())
		}
		slag.Ok("%s exists\n", sysWaresBin)

		path := os.Getenv("PATH")
		if !strings.Contains(path, waresBin) {
			fmt.Print(slag.Err("%s not in PATH\n", waresBin).Error())
		}
		slag.Ok("%s is in PATH\n", waresBin)
		if !strings.Contains(path, sysWaresBin) {
			fmt.Print(slag.Err("%s not in PATH\n", sysWaresBin).Error())
		}
		slag.Ok("%s is in PATH\n", sysWaresBin)

		// MARK: Check for Wares config directory
		waresConfig, _ := internal.ConfigDir()
		if _, err := os.Stat(waresConfig); os.IsNotExist(err) {
			slag.Warn("%s does not exist, creating\n", waresConfig)
		}
		err = os.MkdirAll(waresConfig, 0o755)
		if err != nil {
			fmt.Print(slag.Err("Failed to create wares config dir: %s\n", err).Error())
		}
		slag.Ok("%s exists\n", waresConfig)

		// MARK: Check for authenticated gh CLI
		if status := checkGhCli(); status == true {
			slag.Ok("Logged into GitHub CLI\n")
		} else {
			fmt.Print(slag.Err("Not logged into GitHub CLI\n").Error())
		}

		// MARK: Check for tar and unzip commands
		if checkCommandExists("tar") {
			slag.Ok("tar command found\n")
		} else {
			fmt.Print(slag.Err("tar command not found\n").Error())
		}

		if checkCommandExists("unzip") {
			slag.Ok("unzip command found\n")
		} else {
			fmt.Print(slag.Err("unzip command not found\n").Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
