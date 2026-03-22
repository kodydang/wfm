package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Install kd-wfm binary to ~/.kd-wfm/ and register shell alias",
	Long: `init copies the kd-wfm binary to ~/.kd-wfm/kd-wfm and appends
alias kd-wfm='~/.kd-wfm/kd-wfm' to ~/.zshrc. Safe to re-run.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		home := os.Getenv("HOME")
		installDir := filepath.Join(home, ".kd-wfm")
		installPath := filepath.Join(installDir, "kd-wfm")
		zshrc := filepath.Join(home, ".zshrc")
		aliasLine := fmt.Sprintf("alias kd-wfm='%s'", installPath)

		// --- Binary copy ---
		if err := installBinary(installDir, installPath); err != nil {
			return err
		}

		// --- Alias registration ---
		if err := appendAliasIfMissing(zshrc, aliasLine); err != nil {
			return err
		}

		return nil
	},
}

func installBinary(installDir, installPath string) error {
	src, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}
	src, err = filepath.EvalSymlinks(src)
	if err != nil {
		return fmt.Errorf("could not resolve executable path: %w", err)
	}
	dst, _ := filepath.EvalSymlinks(installPath) // ignore error — dst may not exist yet

	if src == dst || src == installPath {
		fmt.Printf("binary already installed at %s, skipping\n", installPath)
		return nil
	}

	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("could not create %s: %w", installDir, err)
	}

	if err := copyFile(src, installPath, 0755); err != nil {
		return fmt.Errorf("could not copy binary: %w", err)
	}
	fmt.Printf("installed binary to %s\n", installPath)
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func appendAliasIfMissing(zshrc, aliasLine string) error {
	out, err := exec.Command("grep", "-qF", aliasLine, zshrc).CombinedOutput()
	_ = out
	if err == nil {
		fmt.Printf("alias already in %s, skipping\n", zshrc)
		return nil
	}

	f, err := os.OpenFile(zshrc, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", zshrc, err)
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, aliasLine); err != nil {
		return fmt.Errorf("could not write to %s: %w", zshrc, err)
	}
	fmt.Printf("appended alias to %s\n", zshrc)
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
