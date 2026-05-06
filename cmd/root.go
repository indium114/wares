package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var version = "0.5.1"

var rootCmd = &cobra.Command{
	Use:   "wares",
	Short: "A declarative AppImage/binary package manager",
}

func Execute() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(`{{.Version}}`)
	err := fang.Execute(context.Background(), rootCmd)
	if err != nil {
		os.Exit(1)
	}
}
