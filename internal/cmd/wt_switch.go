package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var wtSwitchCmd = &cobra.Command{
	Use:   "wt-switch <branch>",
	Short: "Print the path to a worktree (use via: cd $(kd-wfm wt-switch <branch>))",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]

		container, err := repoContainer()
		if err != nil {
			return err
		}

		dest := filepath.Join(container, filepath.FromSlash(branch))

		if _, err := os.Stat(dest); err != nil {
			return fmt.Errorf("worktree for branch %q not found at %s", branch, dest)
		}

		fmt.Print(dest)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(wtSwitchCmd)
}
