package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ensureBlueprintRepo(repo string) (string, error) {
	home, _ := os.UserHomeDir()
	base := filepath.Join(home, ".local", "share", "wares")

	parts := strings.Split(repo, "/")
	dir := filepath.Join(base, parts[0], parts[1])

	// ensure that repo exists (like the function name :P)
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		// pull latest changes
		fmt.Printf("%s Pulling %s\n", HintText, repo)
		cmd := exec.Command("git", "-C", dir, "pull", "origin")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}
		return dir, nil
	}

	// clone if it doesn't exist
	os.MkdirAll(filepath.Dir(dir), 0o755)
	fmt.Printf("%s Cloning %s\n", HintText, repo)
	cmd := exec.Command("git", "clone", repo, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return dir, nil
}

func resolveLatestCommit(repoDir string) (string, error) {
	cmd := exec.Command("git", "-C", repoDir, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func buildBlueprint(repoDir, commit string, steps []string) error {
	// checkout the locked commit
	cmd := exec.Command("git", "-C", repoDir, "checkout", commit)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// build the project according to steps
	for _, step := range steps {
		fmt.Printf("%s Build step: %s\n", LogText, step)
		cmd := exec.Command("sh", "-c", step)
		cmd.Dir = repoDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
