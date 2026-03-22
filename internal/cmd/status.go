package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show workflow status",
	Long:  `Display the current status of the workflow manager and active workflows.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kd-wfm status: no workflows configured yet")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
