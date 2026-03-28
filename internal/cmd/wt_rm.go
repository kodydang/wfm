package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var wtRmCmd = &cobra.Command{
	Use:   "wt-rm <branch>",
	Short: "Remove a worktree and delete its branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]

		container, err := repoContainer()
		if err != nil {
			return err
		}

		dest := filepath.Join(container, filepath.FromSlash(branch))

		// Remove the worktree (--force handles unmerged changes)
		rm := exec.Command("git", "worktree", "remove", "--force", dest)
		rm.Stdout = os.Stdout
		rm.Stderr = os.Stderr
		if err := rm.Run(); err != nil {
			return fmt.Errorf("git worktree remove failed: %w", err)
		}

		// Delete the branch
		del := exec.Command("git", "branch", "-D", branch)
		del.Stdout = os.Stdout
		del.Stderr = os.Stderr
		if err := del.Run(); err != nil {
			return fmt.Errorf("git branch -D failed: %w", err)
		}

		// Remove empty parent directories up to the container
		if err := removeEmptyParents(dest, container); err != nil {
			// non-fatal — directory may already be gone
			fmt.Fprintf(os.Stderr, "note: could not clean up empty dirs: %v\n", err)
		}

		fmt.Printf("removed worktree %s and branch %s\n", dest, branch)
		return nil
	},
}

// removeEmptyParents removes empty directories from path up to (but not including) stopAt.
func removeEmptyParents(path, stopAt string) error {
	dir := filepath.Dir(path)
	for dir != stopAt && len(dir) > len(stopAt) {
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			break
		}
		if err := os.Remove(dir); err != nil {
			return err
		}
		dir = filepath.Dir(dir)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(wtRmCmd)
}
