# wfm

Kody's Workflow Manager — a personal CLI tool for managing workflow machine operations and automating workflow tasks.

## Installation

```bash
go install github.com/kodydang/kd-wfm@latest
```

Or build from source:

```bash
git clone <repo>
cd kd-wfm
make build
```

## Usage

```bash
# Show help
kd-wfm --help

# Show version
kd-wfm --version

# Show workflow status
kd-wfm status
```

## Development

```bash
# Install dependencies
go mod tidy

# Build
make build

# Run tests
make test

# Install to $GOPATH/bin
make install

# Clean built binary
make clean
```
