package internal

import (
	"os"
	"path/filepath"

	"github.com/indium114/slag"
)

func BuildFromWaresfile(dir string) error {
	home, _ := os.UserHomeDir()

	oldHome := os.Getenv("WARES_HOME")
	os.Setenv("WARES_HOME", filepath.Join(home, ".local", "share", "wares", "_builds"))
	defer os.Setenv("WARES_HOME", oldHome)

	shellConfig, err := LoadShellConfig(dir)
	if err != nil {
		return err
	}

	bp, ok := shellConfig.Blueprints["_self"]
	if !ok {
		return slag.Err("no '_self' blueprint in the specified directory's waresfile.yaml")
	}

	bp.Repo = dir

	slag.Build("%s\n", bp.Repo)

	commit, err := resolveLatestCommit(dir)
	if err != nil {
		return err
	}

	repoDir, err := ensureBlueprintRepo(bp.Repo)
	if err != nil {
		return err
	}

	if err := buildBlueprint(repoDir, commit, bp.Steps); err != nil {
		return err
	}

	for _, artifact := range bp.Artifacts {
		slag.Log("symlinking artifact %s from repoDir %s to result dir %s", artifact, repoDir, dir+"/wares-result") // DEBUG
		if err := shellSymlinkBlueprint(artifact, repoDir, dir+"/wares-result"); err != nil {
			return err
		}
	}

	return nil
}
