## 1. Move Entry Point

- [x] 1.1 Copy `cmd/kd-wfm/main.go` to `main.go` at the repository root
- [x] 1.2 Delete `cmd/kd-wfm/main.go` and remove the `cmd/kd-wfm/` directory
- [x] 1.3 Verify `go build -o kd-wfm .` succeeds from repo root

## 2. Add Makefile

- [x] 2.1 Create `Makefile` at repo root with `build`, `install`, `clean`, and `test` targets
- [x] 2.2 Ensure `make build` outputs binary named `kd-wfm`
- [x] 2.3 Verify `make test` runs `go test ./...` successfully
- [x] 2.4 Add `kd-wfm` binary to `.gitignore`

## 3. Cleanup

- [x] 3.1 Update README build instructions to reference `make build` instead of `go build ./cmd/kd-wfm`
- [x] 3.2 Check for any other references to `cmd/kd-wfm` in the repo and update them
