# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make build      # Build binary to bin/sac with version info via ldflags
make install    # Install globally via go install
make test       # Run all tests
make lint       # Run golangci-lint
```

## Architecture

**sac** is a CLI tool for scaffolding Go projects. Built with Cobra framework.

### Package Structure

- `cmd/sac/main.go` - Entry point, calls `cmd.Execute()`
- `internal/cmd/` - Cobra command definitions
  - `root.go` - Root command with global flags (--debug, --verbose), signal handling
  - `init.go` - `sac project init` command for full project scaffolding
  - `grpc.go` - `sac project grpc` command for gRPC service scaffolding
- `internal/scaffold/` - Project generation logic
  - `scaffold.go` - `Project` struct with templates for full project init
  - `grpc.go` - `GRPCService` struct with proto/gRPC templates
- `internal/config/` - Runtime config (debug/verbose flags)
- `internal/version/` - Build-time version info (set via ldflags)

### Adding New Commands

1. Create file in `internal/cmd/` following existing patterns
2. Define cobra.Command with Use, Short, Long, RunE
3. Register via `init()` with `rootCmd.AddCommand()` or parent command
4. For scaffolding commands, add corresponding generator in `internal/scaffold/`

### Template Pattern

Scaffold generators use a steps pattern:
```go
steps := []struct {
    name string
    fn   func() error
}{
    {"Creating directories", g.createDirs},
    {"Generating files", g.createFiles},
    ...
}
```

Templates returning strings with backticks - escape `%s` as `%%s` when used inside `fmt.Sprintf` raw strings. For struct tags with backticks (e.g., yaml tags), use escaped double-quote strings instead of raw strings.

## CLI Commands

```bash
sac project init <name> [-m module] [-o output]   # Scaffold full project
sac project grpc <service> [-m module] [-o output] # Scaffold gRPC service
sac version                                        # Show version info
```
