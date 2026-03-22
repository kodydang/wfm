## Why

The `main.go` entry point is nested under `cmd/kd-wfm/`, which adds unnecessary indirection for a single-binary CLI. Moving it to the repository root simplifies the project layout and makes building more intuitive. A Makefile standardizes the build process with common targets.

## What Changes

- Move `cmd/kd-wfm/main.go` to the repository root (`main.go`)
- Remove the now-empty `cmd/kd-wfm/` directory
- Add a `Makefile` with targets: `build`, `install`, `clean`, and `test`
- Update any references to the old entry point path (e.g., `go build ./cmd/kd-wfm`)

## Capabilities

### New Capabilities
- `makefile-build`: Standardized Makefile providing `build`, `install`, `clean`, and `test` targets for the CLI binary

### Modified Capabilities
<!-- No spec-level requirement changes — this is a structural/build change only -->

## Impact

- `cmd/kd-wfm/main.go` → `main.go` (file move)
- `cmd/kd-wfm/` directory removed
- `go.mod` module path unchanged; only the build invocation path changes
- `Makefile` added at repository root
- README or CI scripts referencing `go build ./cmd/kd-wfm` will need updating
