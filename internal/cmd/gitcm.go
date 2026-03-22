package cmd

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed prompts/gitcm-system.md
var gitcmSystemPrompt string

var gitcmCmd = &cobra.Command{
	Use:   "gitcm",
	Short: "Generate and apply a conventional git commit message using Claude",
	Long: `gitcm generates a git commit message from your staged changes using the Claude CLI,
shows you the message for confirmation, then commits.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check staged diff
		diff, err := stagedDiff()
		if err != nil {
			return err
		}
		if strings.TrimSpace(diff) == "" {
			return errors.New("no staged changes to commit")
		}

		// Generate commit message via claude CLI
		message, err := generateCommitMessage(diff)
		if err != nil {
			return err
		}

		message = strings.TrimSpace(message)
		if message == "" {
			return errors.New("claude returned an empty commit message")
		}

		// Confirmation prompt
		fmt.Println("\nGenerated commit message:")
		fmt.Println("─────────────────────────")
		fmt.Println(message)
		fmt.Println("─────────────────────────")
		fmt.Print("\nCommit with this message? [y/N] ")

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		if strings.ToLower(answer) != "y" {
			fmt.Println("\nAborted. Generated message:")
			fmt.Println(message)
			return nil
		}

		// Commit
		out, err := runGitCommit(message)
		if err != nil {
			return fmt.Errorf("git commit failed: %w", err)
		}
		fmt.Print(out)
		return nil
	},
}

func stagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff --cached failed: %w", err)
	}
	return out.String(), nil
}

func generateCommitMessage(diff string) (string, error) {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return "", errors.New("claude CLI not found — install Claude Code to use this command (https://claude.ai/code)")
	}

	cmd := exec.Command(claudePath, "-p", gitcmSystemPrompt)
	cmd.Stdin = strings.NewReader(diff)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("claude CLI failed: %w", err)
	}
	return string(out), nil
}

func runGitCommit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func init() {
	rootCmd.AddCommand(gitcmCmd)
}
