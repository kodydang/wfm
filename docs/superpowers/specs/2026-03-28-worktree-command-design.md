# Worktree Command Design

**Date:** 2026-03-28
**Status:** Approved

## Overview

Add four new top-level Cobra commands (`wt-clone`, `wt-add`, `wt-switch`, `wt-rm`) to `kd-wfm` that manage git worktrees with a structured folder layout. The `alias` command is extended to emit a `gwt` shell function into `~/.zshrc` that routes subcommands and handles the `cd` use case for `switch`.

---

## Folder Layout

The repo container is one level above the `main/` worktree. Branch name path segments map directly to the folder path under the repo container.

```
<repo-name>/
  main/               ← main worktree (created by wt-clone)
  feat/
    abc/              ← branch feat/abc (created by gwt add feat/abc)
  hotfix/
    abc/              ← branch hotfix/abc (created by gwt add hotfix/abc)
```

---

## Commands

### `kd-wfm wt-clone <url>`

1. Derives `<repo-name>` from the URL (strip `.git` suffix, take last path segment)
2. Creates directory `<repo-name>/main/`
3. Runs `git clone <url> <repo-name>/main`
4. Prints the path to `<repo-name>/main`

### `kd-wfm wt-add -b <branch>`

Flags:
- `-b <branch>` — branch name to create (required)

Steps:
1. Finds repo root: runs `git rev-parse --show-toplevel` from cwd to get the current worktree path, then walks up one level to get the repo container
2. Derives destination path: `<repo-container>/<branch>` (branch slash becomes path separator)
3. Runs `git worktree add -b <branch> <dest>`
4. Prints the destination path

### `kd-wfm wt-switch <branch>`

1. Finds repo container (same resolution as `wt-add`)
2. Computes path: `<repo-container>/<branch>`
3. Verifies the path exists; errors if not
4. Prints the absolute path to stdout (used by the `gwt switch` shell function via command substitution)

### `kd-wfm wt-rm <branch>`

1. Finds repo container
2. Computes path: `<repo-container>/<branch>`
3. Runs `git worktree remove --force <dest>`
4. Runs `git branch -D <branch>`
5. Removes any leftover empty parent directories up to the repo container

---

## Repo Root Resolution

Used by `wt-add`, `wt-switch`, and `wt-rm`:

1. Run `git worktree list --porcelain` from cwd — the first entry is always the main worktree
2. Extract the main worktree path (e.g. `/projects/my-app/main`)
3. Walk up one level from main worktree path → `/projects/my-app` (the repo container)

This is reliable regardless of how deep the current worktree is (`feat/abc` is 2 levels, `main` is 1 level — the main worktree path is always the anchor).

**Edge case:** if cwd is not inside a git repo, `git worktree list` fails and the command returns a clear error asking the user to `cd` into a worktree first.

---

## Shell Function (`gwt`)

The `alias` command emits this function block into `~/.zshrc`. It is idempotent: the existing `existingKdAliases` check is extended to detect the function by its unique header line `gwt()`.

```bash
gwt() {
  case "$1" in
    clone)  kd-wfm wt-clone "${@:2}" ;;
    add)    kd-wfm wt-add -b "${@:2}" ;;
    switch) cd "$(kd-wfm wt-switch "${@:2}")" ;;
    rm)     kd-wfm wt-rm "${@:2}" ;;
    *)      echo "gwt: unknown command '$1'" >&2; return 1 ;;
  esac
}
```

- `gwt add feat/abc` → injects `-b` automatically
- `gwt switch feat/abc` → `cd` happens in the shell via command substitution
- Checked for presence by scanning `~/.zshrc` for the line `gwt() {` before appending

---

## Changes to Existing Code

### `internal/cmd/alias.go`

- After emitting per-command aliases, also emit the `gwt()` function block
- Check idempotency by scanning for `gwt() {` header line in `~/.zshrc` (using the existing `grep` pattern)
- Append the full multi-line function block as a single write if missing

---

## File Structure

```
internal/cmd/
  wt_clone.go    ← wt-clone command
  wt_add.go      ← wt-add command
  wt_switch.go   ← wt-switch command
  wt_rm.go       ← wt-rm command
  wt_util.go     ← shared repo root resolution logic
  alias.go       ← extended to emit gwt() function
```

---

## Error Handling

| Scenario | Error |
|---|---|
| `wt-clone` URL has no recognizable repo name | `"could not derive repo name from URL"` |
| `wt-add` run outside a git worktree | `"not inside a git worktree — cd into a worktree first"` |
| `wt-switch` target path does not exist | `"worktree for branch '<branch>' not found at <path>"` |
| `wt-rm` branch has unmerged changes | `git worktree remove --force` handles it; branch deleted with `git branch -D` |
