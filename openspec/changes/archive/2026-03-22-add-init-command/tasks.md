## 1. init Command

- [x] 1.1 Create `internal/cmd/init.go` with a Cobra command struct (`Use: "init"`, short description)
- [x] 1.2 Implement binary self-copy: use `os.Executable()` to get source path, `os.MkdirAll` for `~/.kd-wfm/`, copy bytes to `~/.kd-wfm/kd-wfm` with mode `0755`
- [x] 1.3 Detect source == destination: compare resolved paths and skip copy with message if equal
- [x] 1.4 Implement alias registration: use `grep -F` to check if `alias kd-wfm=` already in `~/.zshrc`, append if absent, skip with message if present
- [x] 1.5 Register `init` command in `rootCmd` via `init()` func

## 2. Verification

- [x] 2.1 Build the binary (`go build ./...`) and verify no compile errors
- [x] 2.2 Run `kd-wfm init` and confirm binary appears at `~/.kd-wfm/kd-wfm` and alias is appended to `~/.zshrc`
- [x] 2.3 Re-run `kd-wfm init` and confirm both steps are skipped with informational messages
- [x] 2.4 Run `kd-wfm --help` and confirm `init` appears in the command list
