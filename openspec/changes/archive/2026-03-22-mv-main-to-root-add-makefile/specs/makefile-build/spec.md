## ADDED Requirements

### Requirement: Root entry point
The repository SHALL have `main.go` at the root as the sole entry point for the `kd-wfm` binary. The `cmd/kd-wfm/` directory SHALL be removed.

#### Scenario: Build from root
- **WHEN** a developer runs `go build .` from the repository root
- **THEN** a binary named `kd-wfm` is produced at the repository root

#### Scenario: No cmd subdirectory
- **WHEN** the repository is checked out
- **THEN** `cmd/kd-wfm/` SHALL NOT exist

### Requirement: Makefile build target
The repository SHALL include a `Makefile` at the root with a `build` target that compiles the CLI binary named `kd-wfm`.

#### Scenario: make build produces binary
- **WHEN** a developer runs `make build`
- **THEN** a binary named `kd-wfm` is produced at the repository root

### Requirement: Makefile install target
The `Makefile` SHALL provide an `install` target that installs the binary to the Go bin path.

#### Scenario: make install
- **WHEN** a developer runs `make install`
- **THEN** `kd-wfm` is installed via `go install .` and available in `$GOPATH/bin`

### Requirement: Makefile clean target
The `Makefile` SHALL provide a `clean` target that removes the built binary.

#### Scenario: make clean removes binary
- **WHEN** `make clean` is run after `make build`
- **THEN** the `kd-wfm` binary at the repository root is deleted

### Requirement: Makefile test target
The `Makefile` SHALL provide a `test` target that runs all Go tests.

#### Scenario: make test runs tests
- **WHEN** a developer runs `make test`
- **THEN** `go test ./...` is executed across all packages
