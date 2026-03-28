package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var wtAddBranch string

var wtAddCmd = &cobra.Command{
	Use:   "wt-add",
	Short: "Add a git worktree at <repo-container>/<branch>",
	Long: `Add a new git worktree with the given branch name.
The worktree is created at <repo-container>/<branch>, mirroring
the branch name as a directory path.

Example:
  kd-wfm wt-add -b feat/abc   →  my-app/feat/abc/`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if wtAddBranch == "" {
			return fmt.Errorf("flag -b is required")
		}

		container, err := repoContainer()
		if err != nil {
			return err
		}

		dest := filepath.Join(container, filepath.FromSlash(wtAddBranch))

		c := exec.Command("git", "worktree", "add", "-b", wtAddBranch, dest)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("git worktree add failed: %w", err)
		}

		fmt.Println(dest)
		return nil
	},
}

func init() {
	wtAddCmd.Flags().StringVarP(&wtAddBranch, "branch", "b", "", "Branch name to create (required)")
	rootCmd.AddCommand(wtAddCmd)
}
