package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Print shell alias definitions for all kd-wfm commands",
	Long: `Print shell alias definitions for all kd-wfm subcommands to stdout.

Add to your shell config for dotfiles portability:
  eval "$(kd-wfm alias)"
or:
  kd-wfm alias >> ~/.zshrc`,
	Run: func(cmd *cobra.Command, args []string) {
		skip := map[string]bool{
			"help":       true,
			"completion": true,
			"alias":      true,
		}
		for _, c := range rootCmd.Commands() {
			name := c.Name()
			if skip[name] {
				continue
			}
			fmt.Printf("alias %s='kd-wfm %s'\n", name, name)
		}
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
