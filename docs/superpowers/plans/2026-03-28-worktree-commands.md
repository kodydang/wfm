# Worktree Commands Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `wt-clone`, `wt-add`, `wt-switch`, `wt-rm` commands to `kd-wfm` that manage git worktrees in a structured folder layout, plus a `gwt` shell function emitted by the `alias` command.

**Architecture:** Each command is a focused Cobra command registered via `init()` in its own file, following existing patterns in `internal/cmd/`. Shared path resolution logic lives in `wt_util.go`. The `alias.go` file is extended to emit the `gwt()` shell function idempotently.

**Tech Stack:** Go 1.24, github.com/spf13/cobra, standard library (`os/exec`, `path/filepath`, `strings`)

---

## File Map

| File | Action | Responsibility |
|---|---|---|
| `internal/cmd/wt_util.go` | Create | `parseMainWorktreePath()`, `repoContainer()` — shared by all wt-* commands |
| `internal/cmd/wt_util_test.go` | Create | Unit tests for pure helper functions |
| `internal/cmd/wt_clone.go` | Create | `wt-clone` command + `repoNameFromURL()` |
| `internal/cmd/wt_clone_test.go` | Create | Unit tests for `repoNameFromURL()` |
| `internal/cmd/wt_add.go` | Create | `wt-add` command |
| `internal/cmd/wt_switch.go` | Create | `wt-switch` command |
| `internal/cmd/wt_rm.go` | Create | `wt-rm` command |
| `internal/cmd/alias.go` | Modify | Extend to emit `gwt()` shell function |

---

## Task 1: Shared utilities (`wt_util.go`)

**Files:**
- Create: `internal/cmd/wt_util.go`
- Create: `internal/cmd/wt_util_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/cmd/wt_util_test.go`:

```go
package cmd

import (
	"testing"
)

func TestParseMainWorktreePath(t *testing.T) {
	porcelain := `worktree /projects/my-app/main
HEAD abc123
branch refs/heads/main

worktree /projects/my-app/feat/abc
HEAD def456
branch refs/heads/feat/abc

`
	got, err := parseMainWorktreePath(porcelain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/projects/my-app/main" {
		t.Errorf("got %q, want %q", got, "/projects/my-app/main")
	}
}

func TestParseMainWorktreePath_Empty(t *testing.T) {
	_, err := parseMainWorktreePath("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestRepoContainerFromMainPath(t *testing.T) {
	got := repoContainerFromMainPath("/projects/my-app/main")
	if got != "/projects/my-app" {
		t.Errorf("got %q, want %q", got, "/projects/my-app")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/kodydang/kd-codes/kd-wfm
go test ./internal/cmd/... -run "TestParseMainWorktreePath|TestRepoContainerFromMainPath" -v
```

Expected: FAIL with `undefined: parseMainWorktreePath`

- [ ] **Step 3: Write the implementation**

Create `internal/cmd/wt_util.go`:

```go
package cmd

import (
	"errors"
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
		return "", errors.New("not inside a git worktree — cd into a worktree first")
	}
	mainPath, err := parseMainWorktreePath(string(out))
	if err != nil {
		return "", err
	}
	return repoContainerFromMainPath(mainPath), nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/cmd/... -run "TestParseMainWorktreePath|TestRepoContainerFromMainPath" -v
```

Expected: PASS — 3 tests

- [ ] **Step 5: Commit**

```bash
git add internal/cmd/wt_util.go internal/cmd/wt_util_test.go
git commit -m "feat(wt): add shared worktree utilities"
```

---

## Task 2: `wt-clone` command

**Files:**
- Create: `internal/cmd/wt_clone.go`
- Create: `internal/cmd/wt_clone_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/cmd/wt_clone_test.go`:

```go
package cmd

import "testing"

func TestRepoNameFromURL(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://github.com/user/my-app.git", "my-app"},
		{"https://github.com/user/my-app", "my-app"},
		{"git@github.com:user/my-app.git", "my-app"},
		{"git@github.com:user/my-app", "my-app"},
		{"ssh://git@github.com/user/my-app.git", "my-app"},
	}
	for _, c := range cases {
		got, err := repoNameFromURL(c.input)
		if err != nil {
			t.Errorf("repoNameFromURL(%q): unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("repoNameFromURL(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestRepoNameFromURL_Invalid(t *testing.T) {
	_, err := repoNameFromURL("https://github.com/user/")
	if err == nil {
		t.Fatal("expected error for trailing slash URL, got nil")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/cmd/... -run "TestRepoNameFromURL" -v
```

Expected: FAIL with `undefined: repoNameFromURL`

- [ ] **Step 3: Write the implementation**

Create `internal/cmd/wt_clone.go`:

```go
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
	}
	parts := strings.Split(strings.TrimRight(rawURL, "/"), "/")
	name := parts[len(parts)-1]
	if name == "" {
		return "", errors.New("could not derive repo name from URL")
	}
	return name, nil
}

func init() {
	rootCmd.AddCommand(wtCloneCmd)
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/cmd/... -run "TestRepoNameFromURL" -v
```

Expected: PASS — 6 cases

- [ ] **Step 5: Build and smoke-test**

```bash
go build -o /tmp/kd-wfm . && /tmp/kd-wfm wt-clone --help
```

Expected: help text showing `wt-clone <url>`

- [ ] **Step 6: Commit**

```bash
git add internal/cmd/wt_clone.go internal/cmd/wt_clone_test.go
git commit -m "feat(wt): add wt-clone command"
```

---

## Task 3: `wt-add` command

**Files:**
- Create: `internal/cmd/wt_add.go`

- [ ] **Step 1: Write the implementation**

Create `internal/cmd/wt_add.go`:

```go
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
```

- [ ] **Step 2: Build and verify**

```bash
go build -o /tmp/kd-wfm . && /tmp/kd-wfm wt-add --help
```

Expected: help text showing `-b` flag

- [ ] **Step 3: Commit**

```bash
git add internal/cmd/wt_add.go
git commit -m "feat(wt): add wt-add command"
```

---

## Task 4: `wt-switch` command

**Files:**
- Create: `internal/cmd/wt_switch.go`

- [ ] **Step 1: Write the implementation**

Create `internal/cmd/wt_switch.go`:

```go
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
```

- [ ] **Step 2: Build and verify**

```bash
go build -o /tmp/kd-wfm . && /tmp/kd-wfm wt-switch --help
```

Expected: help text showing `wt-switch <branch>`

- [ ] **Step 3: Commit**

```bash
git add internal/cmd/wt_switch.go
git commit -m "feat(wt): add wt-switch command"
```

---

## Task 5: `wt-rm` command

**Files:**
- Create: `internal/cmd/wt_rm.go`

- [ ] **Step 1: Write the implementation**

Create `internal/cmd/wt_rm.go`:

```go
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
```

- [ ] **Step 2: Build and verify**

```bash
go build -o /tmp/kd-wfm . && /tmp/kd-wfm wt-rm --help
```

Expected: help text showing `wt-rm <branch>`

- [ ] **Step 3: Commit**

```bash
git add internal/cmd/wt_rm.go
git commit -m "feat(wt): add wt-rm command"
```

---

## Task 6: Extend `alias.go` to emit `gwt()` shell function

**Files:**
- Modify: `internal/cmd/alias.go`

The current `alias.go` uses `existingKdAliases()` to grep for `kd-wfm` lines and skip duplicates. We extend it to also emit the `gwt()` function block, checking idempotency by looking for the header line `gwt() {`.

- [ ] **Step 1: Read the current `alias.go`**

File is at `internal/cmd/alias.go`. Key function to extend: the `RunE` body of `aliasCmd`. After the existing per-command alias loop, add the `gwt()` function block.

- [ ] **Step 2: Write the updated `alias.go`**

Replace the entire `RunE` body (lines 17–59) with the extended version:

```go
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

		// gwt() shell function
		const gwtHeader = "gwt() {"
		const gwtFunction = `gwt() {
  case "$1" in
    clone)  kd-wfm wt-clone "${@:2}" ;;
    add)    kd-wfm wt-add -b "${@:2}" ;;
    switch) cd "$(kd-wfm wt-switch "${@:2}")" ;;
    rm)     kd-wfm wt-rm "${@:2}" ;;
    *)      echo "gwt: unknown command '$1'" >&2; return 1 ;;
  esac
}`
		fmt.Println(gwtHeader)
		if gwtFunctionPresent(zshrc, gwtHeader) {
			fmt.Printf("  already in %s, skipping\n", zshrc)
		} else {
			toAppend.WriteString(gwtFunction + "\n")
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
```

Also add the `gwtFunctionPresent` helper after `existingKdAliases`:

```go
// gwtFunctionPresent checks if the gwt() function header line exists in path.
func gwtFunctionPresent(path, header string) bool {
	out, err := exec.Command("grep", "-qF", header, path).CombinedOutput()
	_ = out
	return err == nil
}
```

- [ ] **Step 3: Build and verify**

```bash
go build -o /tmp/kd-wfm . && /tmp/kd-wfm alias --help
```

Expected: builds cleanly, help text shows `alias` command

- [ ] **Step 4: Verify all new commands appear in help**

```bash
/tmp/kd-wfm --help
```

Expected output includes: `wt-add`, `wt-clone`, `wt-rm`, `wt-switch`

- [ ] **Step 5: Run all tests**

```bash
go test ./internal/cmd/... -v
```

Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/cmd/alias.go
git commit -m "feat(wt): extend alias command to emit gwt() shell function"
```

---

## Task 7: End-to-end smoke test

This task verifies the full workflow manually. No code changes.

- [ ] **Step 1: Install the new binary**

```bash
go build -o /tmp/kd-wfm . && /tmp/kd-wfm --help
```

Expected: all 8+ commands listed including `wt-add`, `wt-clone`, `wt-rm`, `wt-switch`

- [ ] **Step 2: Test `wt-clone`**

```bash
cd /tmp
/tmp/kd-wfm wt-clone https://github.com/kodydang/kd-wfm.git
ls /tmp/kd-wfm-repo/main   # should contain cloned repo
```

Expected: directory `/tmp/kd-wfm/main/` created with cloned content

- [ ] **Step 3: Test `wt-add` and `wt-switch`**

```bash
cd /tmp/kd-wfm/main
/tmp/kd-wfm wt-add -b feat/test-branch
ls /tmp/kd-wfm/feat/test-branch   # should exist
/tmp/kd-wfm wt-switch feat/test-branch  # should print path
```

Expected: path printed, directory exists

- [ ] **Step 4: Test `wt-rm`**

```bash
cd /tmp/kd-wfm/main
/tmp/kd-wfm wt-rm feat/test-branch
```

Expected: worktree removed, branch deleted, `feat/` dir cleaned up

- [ ] **Step 5: Test `alias` command emits `gwt()`**

```bash
/tmp/kd-wfm alias
grep "gwt() {" ~/.zshrc
```

Expected: `gwt() {` line present in `~/.zshrc`

- [ ] **Step 6: Final commit (if any fixes applied)**

```bash
git add -p
git commit -m "fix(wt): smoke test fixes"
```
