## Why

Setting up `kd-wfm` on a new machine requires manually placing the binary and wiring up a shell alias. An `init` command automates this one-time setup so the CLI is usable immediately after running a single command.

## What Changes

- Add an `init` subcommand to `kd-wfm` that installs the binary to `~/.kd-wfm/kd-wfm` and appends an alias `kd-wfm` pointing to that path in `~/.zshrc`
- The command is idempotent: re-running it skips steps that are already done

## Capabilities

### New Capabilities
- `init-command`: One-shot setup command that copies the running binary to `~/.kd-wfm/` and registers a shell alias for it

### Modified Capabilities
<!-- None -->

## Impact

- Adds `internal/cmd/init.go`
- Uses `os.Executable()` to locate the running binary at runtime
- Writes to `~/.kd-wfm/` (creates dir if needed) and appends to `~/.zshrc`
- No new Go module dependencies
