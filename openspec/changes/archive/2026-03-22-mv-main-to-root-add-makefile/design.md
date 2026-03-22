## Context

Currently `main.go` lives at `cmd/kd-wfm/main.go`. The file is minimal (3 lines — imports `internal/cmd` and calls `Execute()`). Moving it to the root makes the project layout match single-binary Go CLI conventions and removes a directory level with no structural value.

A `Makefile` will provide standard developer targets so contributors don't need to remember raw `go build` invocations.

## Goals / Non-Goals

**Goals:**
- Move `cmd/kd-wfm/main.go` → `main.go` at repo root
- Delete `cmd/kd-wfm/` directory
- Add `Makefile` with `build`, `install`, `clean`, and `test` targets
- Binary output name: `kd-wfm`

**Non-Goals:**
- Changing the module path or internal package structure
- Cross-platform build matrix (CI/release builds are out of scope here)
- Adding ldflags for version injection (can be a follow-up)

## Decisions

### Move `main.go` to repo root (not keep `cmd/` layout)

The `cmd/<binary>/` pattern is useful when a repo builds multiple binaries. This project has one binary and no plans for more. A flat root `main.go` is idiomatic for single-binary CLIs.

**Alternative considered**: Keep `cmd/kd-wfm/` — rejected because it adds depth with no benefit for a single-binary repo.

### Makefile over shell scripts

A `Makefile` is universally available on macOS/Linux without additional tooling, integrates naturally with editors, and is the de facto standard for Go project build automation.

**Alternative considered**: `Taskfile.yml` — rejected to avoid introducing a dependency (`task` CLI) not already present.

### Makefile targets

| Target | Command | Notes |
|--------|---------|-------|
| `build` | `go build -o kd-wfm .` | Outputs `./kd-wfm` at repo root |
| `install` | `go install .` | Installs to `$GOPATH/bin` |
| `clean` | `rm -f kd-wfm` | Removes built binary |
| `test` | `go test ./...` | Runs all tests |

`.PHONY` declared for all targets.

## Risks / Trade-offs

- **Stale docs/CI** → Any README or CI config referencing `./cmd/kd-wfm` must be updated. Low risk — the repo is early stage.
- **`go build .`** now works from root, which is simpler but changes the prior invocation. Acceptable trade-off given project stage.
