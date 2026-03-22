## Context

`kd-wfm` is a personal CLI built with Cobra/Viper in Go. Currently it has placeholder commands (`status`). This change adds two real workflow commands:

1. **`gitcm`** — generates a conventional commit message by calling the `claude` CLI (Claude Code) with a system prompt stored in the repo, then lets the user confirm before committing
2. **`alias`** — prints shell alias definitions for all `kd-wfm` subcommands to stdout, enabling dotfiles portability across machines

## Goals / Non-Goals

**Goals:**
- `gitcm` calls `claude` CLI subprocess with the repo-stored system prompt and staged diff as input, parses the output as the commit message, shows it for confirmation, then runs `git commit`
- `alias` prints `alias gitcm='kd-wfm gitcm'`-style lines to stdout for all registered commands
- System prompt lives at `prompts/gitcm-system.md` in the repo root, embedded into the binary via `//go:embed`
- No new external Go dependencies beyond what's already in go.mod

**Non-Goals:**
- Auto-staging files (`git add`) — user controls staging
- Appending to shell config files directly
- Supporting LLM APIs other than `claude` CLI (for now)
- Editing or retrying the generated commit message interactively

## Decisions

### Use `claude` CLI subprocess instead of Anthropic Go SDK
`claude` CLI is already installed as part of the user's tooling. Shelling out avoids adding a new Go module dependency and an API key management problem. The `gitcm` command runs `claude -p "<system-prompt>" < diff` and captures stdout as the commit message.

**Alternatives considered:**
- Direct Anthropic HTTP API: requires managing `ANTHROPIC_API_KEY` in the binary, adds `net/http` complexity
- Anthropic Go SDK: adds a dependency and key management; no benefit over using the CLI that's already present

### Embed system prompt via `//go:embed`
The prompt at `prompts/gitcm-system.md` is embedded at compile time. This keeps the binary self-contained and versioned with the code, while staying human-readable and editable in the repo.

**Alternatives considered:**
- Runtime file read from a known path: fragile across machines
- Hardcoded string in Go: hard to edit, clutters source

### Confirmation via simple stdin prompt
After generating the message, print it and ask `Commit with this message? [y/N]`. If denied, print the message and exit 0 so the user can use it manually.

### `alias` command iterates `rootCmd.Commands()`
Rather than maintaining a hardcoded list, `alias` iterates Cobra's registered command list at runtime and prints one `alias <name>='kd-wfm <name>'` line per command (excluding `completion`, `help`).

## Risks / Trade-offs

- **`claude` CLI not installed** → `gitcm` fails with a clear error: "claude CLI not found — install Claude Code to use this command". No silent failure.
- **Empty staged diff** → `gitcm` exits early with "no staged changes to commit".
- **LLM output not a valid commit message** → shown to user during confirmation; they can abort and commit manually.
- **`alias` output shell compatibility** → single-quote style aliases work in both bash and zsh; no shell detection needed.

## Migration Plan

No migration needed. Both commands are new additions. Users who want the aliases add `eval "$(kd-wfm alias)"` or `kd-wfm alias >> ~/.zshrc` to their shell config.
