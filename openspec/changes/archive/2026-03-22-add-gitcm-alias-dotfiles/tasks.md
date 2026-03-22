## 1. System Prompt

- [x] 1.1 Create `prompts/gitcm-system.md` with a conventional commit generation prompt (instruct the model to output only the commit message, no explanation)

## 2. gitcm Command

- [x] 2.1 Create `internal/cmd/gitcm.go` with a Cobra command struct (`Use: "gitcm"`, short description)
- [x] 2.2 Add `//go:embed` directive to load `prompts/gitcm-system.md` into a package-level variable
- [x] 2.3 Implement staged diff check: run `git diff --cached` and exit early with message if output is empty
- [x] 2.4 Implement `claude` CLI invocation: exec `claude -p <system-prompt>` with diff piped to stdin, capture stdout
- [x] 2.5 Handle `claude` not found in PATH: print actionable error message and exit non-zero
- [x] 2.6 Implement confirmation prompt: print generated message, ask `[y/N]`, read stdin response
- [x] 2.7 On confirm: run `git commit -m "<message>"` and print output; on decline: print message and exit 0
- [x] 2.8 Register `gitcm` command in `internal/cmd/root.go`

## 3. alias Command

- [x] 3.1 Create `internal/cmd/alias.go` with a Cobra command struct (`Use: "alias"`, short description)
- [x] 3.2 Implement alias generation: iterate `rootCmd.Commands()`, skip `help`/`completion`/`alias`, print `alias <name>='kd-wfm <name>'` per command
- [x] 3.3 Register `alias` command in `internal/cmd/root.go`

## 4. Verification

- [x] 4.1 Build the binary (`go build ./...`) and verify no compile errors
- [x] 4.2 Run `kd-wfm alias` and confirm `gitcm` alias appears in output
- [x] 4.3 Stage a file change and run `kd-wfm gitcm` end-to-end to verify the confirmation flow
- [x] 4.4 Run `kd-wfm --help` and confirm both `gitcm` and `alias` appear
