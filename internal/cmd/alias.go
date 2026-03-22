package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Append shell alias definitions for all kd-wfm commands to ~/.zshrc",
	Long:  `Append shell alias definitions for all kd-wfm subcommands to ~/.zshrc.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skip := map[string]bool{
			"help":       true,
			"completion": true,
			"alias":      true,
		}

		zshrc := filepath.Join(os.Getenv("HOME"), ".zshrc")
		existingAliases := existingKdAliases(zshrc)

		var toAppend strings.Builder
		for _, c := range rootCmd.Commands() {
			name := c.Name()
			if skip[name] {
				continue
			}
			line := fmt.Sprintf("alias %s='kd-wfm %s'", name, name)
			fmt.Println(line)
			if existingAliases[line] {
				fmt.Printf("  already in %s, skipping\n", zshrc)
			} else {
				toAppend.WriteString(line + "\n")
			}
		}

		if toAppend.Len() == 0 {
			fmt.Println("\nAll aliases already present — nothing to append.")
			return nil
		}

		f, err := os.OpenFile(zshrc, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("could not open %s: %w", zshrc, err)
		}
		defer f.Close()

		if _, err := fmt.Fprint(f, toAppend.String()); err != nil {
			return fmt.Errorf("could not write to %s: %w", zshrc, err)
		}

		fmt.Printf("\nAppended to %s\n", zshrc)
		return nil
	},
}

// existingKdAliases uses grep to find only kd-wfm alias lines in path,
// so the file content (tokens, secrets, etc.) is never read into Go memory.
func existingKdAliases(path string) map[string]bool {
	found := make(map[string]bool)
	out, err := exec.Command("grep", "-F", "kd-wfm", path).Output()
	if err != nil {
		// exit code 1 means no matches — not an error
		return found
	}
	for _, line := range strings.Split(strings.TrimRight(string(out), "\n"), "\n") {
		if strings.HasPrefix(line, "alias ") {
			found[line] = true
		}
	}
	return found
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
