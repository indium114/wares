package cmd

import (
	"context"
	"fmt"
	"io"
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
	err := fang.Execute(context.Background(), rootCmd, fang.WithVersion(Version), fang.WithErrorHandler(func(w io.Writer, _ fang.Styles, err error) {
		fmt.Fprintln(w, err.Error())
	}))
	if err != nil {
		os.Exit(1)
	}
}
