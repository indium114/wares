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

func linkBlueprintArtifacts(repoDir string, artifacts []string) error {
	home, _ := os.UserHomeDir()
	waresDir := filepath.Join(home, "Wares")

	for _, artifact := range artifacts {
		src := filepath.Join(repoDir, artifact)

		if _, err := os.Stat(src); err != nil {
			fmt.Printf("%s Artifact %s not found\n", ErrText, artifact)
			return err
		}

		linkPath := filepath.Join(waresDir, filepath.Base(artifact))

		os.Remove(linkPath)

		if err := os.Symlink(src, linkPath); err != nil {
			return err
		}
	}

	return nil
}

func findBlueprintOrphans(cfg *Config, lock *Lockfile) []string {
	var orphans []string

	for name := range lock.Blueprints {
		if _, ok := cfg.Blueprints[name]; !ok {
			orphans = append(orphans, name)
		}
	}

	return orphans
}

func uninstallBlueprintOrphans(cfg *Config, lock *Lockfile) (bool, error) {
	orphans := findBlueprintOrphans(cfg, lock)
	changed := false

	for _, name := range orphans {
		fmt.Printf("%s Removing %s\n", SyncText, name)

		locked := lock.Blueprints[name]

		// unlink
		home, _ := os.UserHomeDir()
		waresDir := filepath.Join(home, "Wares")
		for _, artifact := range locked.Artifacts {
			linkPath := filepath.Join(waresDir, filepath.Base(artifact))
			if err := os.Remove(linkPath); err != nil {
				return false, err
			}
		}

		// clean cloned repo
		if locked.Repo != "" {
			home, _ := os.UserHomeDir()
			base := filepath.Join(home, ".local", "share", "wares")

			parts := strings.Split(locked.Repo, "/")
			if len(parts) == 2 {
				repoDir := filepath.Join(base, parts[0], parts[1])
				if err := os.RemoveAll(repoDir); err != nil {
					return false, err
				}
			}
		}

		// unlock
		delete(lock.Blueprints, name)
		changed = true
	}

	return changed, nil
}

func SyncBlueprints(cfg *Config, lock *Lockfile) (bool, error) {
	changed := false

	if lock.Blueprints == nil {
		lock.Blueprints = map[string]LockedBlueprint{}
	}

	for name, bp := range cfg.Blueprints {
		fmt.Printf("%s Building %s\n", SyncText, name)

		// clone
		repoDir, err := ensureBlueprintRepo(bp.Repo)
		if err != nil {
			return false, err
		}

		// get latest commit
		commit, err := resolveLatestCommit(repoDir)
		if err != nil {
			return false, err
		}

		// don't unnecessarily rebuild
		locked := lock.Blueprints[name]
		needRebuild := locked.Commit != commit || locked.Repo != bp.Repo
		if !needRebuild {
			continue
		}

		// build
		if err := buildBlueprint(repoDir, commit, bp.Steps); err != nil {
			return false, err
		}

		// symlink
		if err := linkBlueprintArtifacts(repoDir, bp.Artifacts); err != nil {
			return false, err
		}

		// lock
		lock.Blueprints[name] = LockedBlueprint{
			Repo:      bp.Repo,
			Commit:    commit,
			Artifacts: bp.Artifacts,
		}
		changed = true
	}

	orphanChanged, err := uninstallBlueprintOrphans(cfg, lock)
	if err != nil {
		return false, err
	}

	changed = changed || orphanChanged

	return changed, nil
}
