package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// parseMainWorktreePath parses `git worktree list --porcelain` output
// and returns the path of the first (main) worktree.
func parseMainWorktreePath(porcelain string) (string, error) {
	for _, line := range strings.Split(porcelain, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			return strings.TrimPrefix(line, "worktree "), nil
		}
	}
	return "", errors.New("could not parse main worktree path from git output")
}

// repoContainerFromMainPath returns the repo container dir
// by walking up one level from the main worktree path.
func repoContainerFromMainPath(mainPath string) string {
	return filepath.Dir(mainPath)
}

// repoContainer finds the repo container directory by running
// `git worktree list --porcelain` and walking up from the main worktree.
func repoContainer() (string, error) {
	out, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return "", fmt.Errorf("not inside a git worktree — cd into a worktree first: %w", err)
	}
	mainPath, err := parseMainWorktreePath(string(out))
	if err != nil {
		return "", err
	}
	return repoContainerFromMainPath(mainPath), nil
}
