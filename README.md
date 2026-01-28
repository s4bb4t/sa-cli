# sac

A CLI tool for scaffolding production-ready Go projects with clean architecture.

## Installation

```bash
# From source
go install github.com/s4bb4t/sa-cli/cmd/sac@latest

# Or clone and build
git clone https://github.com/s4bb4t/sa-cli.git
cd sa-cli
make install
```

## Commands

### `sac project init`

Initialize a new Go project with production-ready structure.

```bash
sac project init myapp
sac project init myapp --module github.com/myorg/myapp
sac project init myapp --output ./projects/myapp
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--module` | `-m` | Go module path (default: project name) |
| `--output` | `-o` | Output directory (default: project name) |

**Generated structure:**
```
myapp/
├── cmd/myapp/main.go           # Entry point with graceful shutdown
├── internal/
│   ├── config/
│   │   ├── config.go           # YAML config loader
│   │   └── otel.go             # OpenTelemetry config
│   ├── infrastructure/
│   │   └── database/
│   └── presentation/
├── pkg/grpc/
├── api/
│   ├── proto/
│   └── openapi/
├── deploy/
│   ├── docker/Dockerfile       # Multi-stage Docker build
│   └── k8s/
├── test/
│   ├── integration/
│   └── e2e/
├── scripts/
├── Makefile
├── .gitignore
└── go.mod
```

**Generated project uses:**
- [zapang](https://github.com/s4bb4t/zapang) for structured logging
- [yaml.v3](https://gopkg.in/yaml.v3) for configuration
- Graceful shutdown with signal handling

---

### `sac project grpc`

Generate a gRPC service scaffold with proto contract.

```bash
sac project grpc user
sac project grpc user --module github.com/myorg/user-service
sac project grpc user --output ./services/user
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--module` | `-m` | Go module path (default: service name) |
| `--output` | `-o` | Output directory (default: current dir) |

**Generated structure:**
```
./
├── api/proto/v1/
│   └── <service>.proto         # Proto contract
├── pkg/grpc/<service>/v1/
│   └── *.pb.go                 # Generated protobuf files
├── internal/presentation/grpc/v1/
├── Makefile                    # With proto target
└── go.mod
```

**After generation:**
```bash
make proto    # Regenerate pb files from proto
make build    # Build the service
```

**Requirements:**
- `protoc` (Protocol Buffers compiler)
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins

---

### `sac version`

Print version information.

```bash
sac version
sac version --short
```

**Output:**
```
sa-cli v1.0.0
  Commit:     abc1234
  Built:      2024-01-15T10:30:00Z
  Go version: go1.21.0
  OS/Arch:    darwin/arm64
```

---

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--debug` | `-d` | Enable debug mode |
| `--verbose` | `-v` | Enable verbose output |
| `--help` | `-h` | Help for any command |

---

## Shell Completion

Generate autocompletion scripts for your shell:

```bash
# Bash
sac completion bash > /etc/bash_completion.d/sac

# Zsh
sac completion zsh > "${fpath[1]}/_sac"

# Fish
sac completion fish > ~/.config/fish/completions/sac.fish

# PowerShell
sac completion powershell > sac.ps1
```

---

## Development

```bash
# Build
make build

# Install locally
make install

# Run tests
make test

# Lint
make lint

# Clean build artifacts
make clean
```

---

## License

MIT
