package cmd

import (
	"fmt"
	"os"

	"github.com/kodydang/kd-wfm/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kd-wfm",
	Short: "kd-wfm is a personal workflow manager",
	Long: `kd-wfm (Kody's Workflow Manager) is a CLI tool for managing
personal workflow machine operations and automating workflow tasks.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version.String()

	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
	})
}
