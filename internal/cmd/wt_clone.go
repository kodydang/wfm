package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var wtCloneCmd = &cobra.Command{
	Use:   "wt-clone <url>",
	Short: "Clone a repo into <repo-name>/main/ worktree layout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		repoName, err := repoNameFromURL(url)
		if err != nil {
			return err
		}

		dest := filepath.Join(repoName, "main")
		if err := os.MkdirAll(dest, 0755); err != nil {
			return fmt.Errorf("could not create %s: %w", dest, err)
		}

		c := exec.Command("git", "clone", url, dest)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}

		abs, err := filepath.Abs(dest)
		if err != nil {
			return err
		}
		fmt.Println(abs)
		return nil
	},
}

// repoNameFromURL derives repo name from a git URL.
// https://github.com/user/my-app.git → my-app
// git@github.com:user/my-app.git     → my-app
// ssh://git@github.com/user/my-app   → my-app
func repoNameFromURL(rawURL string) (string, error) {
	rawURL = strings.TrimSuffix(rawURL, ".git")
	// Handle SCP-style SSH: git@github.com:user/repo (no scheme prefix)
	// Exclude scheme URLs like https:// or ssh:// which have "://"
	if !strings.Contains(rawURL, "://") {
		if idx := strings.Index(rawURL, ":"); idx != -1 {
			rawURL = rawURL[idx+1:]
		}
	} else {
		// For scheme URLs (https://, ssh://), extract path component
		if idx := strings.Index(rawURL, "://"); idx != -1 {
			rawURL = rawURL[idx+3:]
			// Skip host (find the first /)
			if idx := strings.Index(rawURL, "/"); idx != -1 {
				rawURL = rawURL[idx+1:]
			}
		}
	}

	trimmed := strings.TrimRight(rawURL, "/")
	if trimmed == "" {
		return "", errors.New("could not derive repo name from URL")
	}
	parts := strings.Split(trimmed, "/")
	// Filter out empty parts
	var nonEmpty []string
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}
	// Need at least user/repo (2 parts) to have a valid repo name
	if len(nonEmpty) < 2 {
		return "", errors.New("could not derive repo name from URL")
	}
	name := nonEmpty[len(nonEmpty)-1]
	return name, nil
}

func init() {
	rootCmd.AddCommand(wtCloneCmd)
}
