## Why

Setting up a new machine requires manually re-adding shell aliases for personal workflow commands, which is tedious and error-prone. A dotfiles-friendly alias generation command solves this by letting `kd-wfm` output shell aliases that can be eval'd or sourced, making the workflow portable across machines.

## What Changes

- Add a `gitcm` subcommand that stages all changes and generates a commit message via an LLM call using a system prompt stored in the repo, then commits
- Add an `alias` subcommand that prints shell alias definitions (e.g., `alias gitcm='kd-wfm gitcm'`) for all workflow commands, suitable for appending to `.zshrc` or `.bashrc`

## Capabilities

### New Capabilities
- `gitcm-command`: LLM-powered git commit command that reads staged/unstaged changes, calls the Claude API with a repo-stored system prompt, and commits with the generated message following conventional commit conventions
- `alias-output`: Command that prints shell-compatible alias definitions for all `kd-wfm` subcommands to stdout, enabling dotfiles portability

### Modified Capabilities
<!-- None -->

## Impact

- Adds `internal/cmd/gitcm.go` and `internal/cmd/alias.go`
- Adds `prompts/gitcm-system.md` (or similar) to store the LLM system prompt in the repo
- Requires Anthropic SDK dependency (`github.com/anthropic-ai/sdk` or HTTP call to Claude API)
- Users source or eval the alias output in their shell config for machine portability
