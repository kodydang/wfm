# kd-wfm

Kody's Workflow Manager — a personal CLI tool for managing workflow machine operations and automating workflow tasks.

## Installation

```bash
go install github.com/kodydang/kd-wfm/cmd/kd-wfm@latest
```

Or build from source:

```bash
git clone <repo>
cd kd-wsm
go build -o kd-wfm ./cmd/kd-wfm
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
go build ./cmd/kd-wfm

# Build with version info
go build -ldflags "-X github.com/kodydang/kd-wfm/internal/version.Version=v0.1.0 \
  -X github.com/kodydang/kd-wfm/internal/version.Commit=$(git rev-parse --short HEAD) \
  -X github.com/kodydang/kd-wfm/internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o kd-wfm ./cmd/kd-wfm
```
