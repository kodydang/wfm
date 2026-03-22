## Context

We're initializing a new Golang CLI project called `kd-wfm` (Kody's Workflow Manager) to manage personal workflow machine operations. Currently, there's no dedicated tooling for workflow management, and this CLI will serve as the foundation for future workflow automation capabilities.

This is a greenfield project with no existing codebase constraints. The project will live in the current repository and follow Go best practices for CLI tool development.

## Goals / Non-Goals

**Goals:**
- Set up a working Golang CLI project with proper module structure
- Establish a scalable command framework that supports subcommands
- Create a clean project structure that's easy to extend
- Provide a functional CLI with help system and version information
- Follow Go community conventions and best practices

**Non-Goals:**
- Implementing actual workflow management logic (that comes later)
- Building a complex plugin system (keep it simple initially)
- Creating a GUI or web interface
- Multi-platform compilation setup (can be added later)

## Decisions

### CLI Framework: Cobra + Viper

**Decision**: Use `spf13/cobra` for command structure and `spf13/viper` for configuration management.

**Rationale**:
- Cobra is the de facto standard for Go CLIs (used by kubectl, hugo, gh, etc.)
- Built-in support for subcommands, flags, and help generation
- Viper integrates seamlessly for config file and environment variable support
- Well-documented and actively maintained

**Alternatives Considered**:
- `urfave/cli`: Simpler but less powerful for complex command hierarchies
- Custom flag parsing: Too much boilerplate for little benefit

### Project Structure

**Decision**: Use standard Go project layout with cmd/tool separation.

```
kd-wfm/
├── cmd/
│   └── kd-wfm/
│       └── main.go          # Entry point
├── internal/
│   ├── cmd/                 # Command implementations
│   │   └── root.go         # Root command
│   └── version/            # Version info
│       └── version.go
├── go.mod
├── go.sum
└── README.md
```

**Rationale**:
- `cmd/kd-wfm/` allows for potential future multi-binary support
- `internal/` prevents external imports of internal packages
- Separates command logic from main entry point for testability

**Alternatives Considered**:
- Flat structure: Doesn't scale well
- More complex layouts (pkg/, api/): Premature for initial setup

### Go Module Configuration

**Decision**: Initialize as `github.com/kodydang/kd-wfm` (or appropriate path based on repository setup).

**Rationale**:
- Standard Go module naming convention
- Allows for future publishing/sharing if needed
- Aligns with repository structure

### Version Management

**Decision**: Use build-time variable injection for version information.

**Rationale**:
- Allows version to be set during build without code changes
- Standard practice for Go CLI tools
- Supports `--version` flag easily

## Risks / Trade-offs

**Risk**: Cobra adds dependency weight for a simple initial CLI
→ **Mitigation**: The initial overhead is justified by long-term maintainability and extensibility. Cobra prevents reinventing command parsing.

**Risk**: Choosing wrong project structure early could require refactoring
→ **Mitigation**: Following standard Go project layout minimizes risk. Structure is well-established in the Go community.

**Trade-off**: Using `internal/` package prevents external reuse
→ **Rationale**: This is a personal tool, not a library. If components need to be shared later, they can be extracted to `pkg/`.
