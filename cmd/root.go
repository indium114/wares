package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "wares",
	Short: "A declarative AppImage/binary package manager",
}

func Execute() {
	err := fang.Execute(context.Background(), rootCmd, fang.WithVersion(Version))
	if err != nil {
		os.Exit(1)
	}
}
