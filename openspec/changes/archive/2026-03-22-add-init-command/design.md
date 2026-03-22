## Context

`kd-wfm` is a personal Go CLI. After building the binary, users must manually copy it somewhere on their PATH and add a shell alias. The `init` command automates this, similar to how other CLI tools provide a self-install step.

## Goals / Non-Goals

**Goals:**
- Copy the running binary to `~/.kd-wfm/kd-wfm` (create directory if needed)
- Append `alias kd-wfm='~/.kd-wfm/kd-wfm'` to `~/.zshrc` if not already present
- Idempotent: re-running prints what is already done and skips those steps

**Non-Goals:**
- Modifying `$PATH` (alias approach avoids this)
- Supporting shells other than zsh (can be extended later)
- Uninstall / removing the binary or alias

## Decisions

### Use `os.Executable()` to locate the source binary
`os.Executable()` returns the path of the running binary, so `init` can copy itself without the user specifying a path.

**Alternatives considered:**
- Accept a `--binary` flag: adds friction; the common case is self-install
- Hard-code `./kd-wfm`: fragile — depends on CWD

### Install to `~/.kd-wfm/kd-wfm`, alias as `kd-wfm`
A dedicated `~/.kd-wfm/` directory keeps the binary isolated and easy to find. The alias in `~/.zshrc` makes it available in new shells without touching `$PATH`.

**Alternatives considered:**
- `~/.local/bin/`: more standard but requires `$PATH` setup, which varies per machine
- `/usr/local/bin/`: requires sudo

### Reuse `existingKdAliases` pattern from `alias.go`
Use `grep -F` to check for the alias line before appending, consistent with the existing alias command's approach.

## Risks / Trade-offs

- **Binary permissions** → Copy with mode `0755` so the installed binary is executable
- **Source == destination** (user already ran from `~/.kd-wfm/kd-wfm`) → Detect and skip the copy step with a message
- **`~/.zshrc` does not exist** → `os.OpenFile` with `O_CREATE` handles this gracefully
