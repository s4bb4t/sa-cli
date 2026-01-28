package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	dirPerm  = 0755
	filePerm = 0644
)

type Project struct {
	Name      string
	Module    string
	OutputDir string
}

func New(name, module, outputDir string) *Project {
	if module == "" {
		module = name
	}
	return &Project{
		Name:      name,
		Module:    module,
		OutputDir: outputDir,
	}
}

func (p *Project) Generate() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Creating directories", p.createDirs},
		{"Generating files", p.createFiles},
		{"Initializing go module", p.initGoMod},
		{"Installing dependencies", p.goModTidy},
	}

	for _, step := range steps {
		fmt.Printf("  →   %s...\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

func (p *Project) createDirs() error {
	dirs := []string{
		// Commands
		"cmd/" + p.Name,

		// API definitions
		"api/proto",
		"api/openapi",

		// Internal packages (private)
		"internal/config",
		"internal/infrastructure/database",
		"internal/presentation",

		// Public packages
		"pkg/grpc",

		// Deployments
		"deploy/docker",
		"deploy/k8s",

		// Scripts
		"scripts",

		// Tests
		"test/integration",
		"test/e2e",
	}

	for _, dir := range dirs {
		path := filepath.Join(p.OutputDir, dir)
		if err := os.MkdirAll(path, dirPerm); err != nil {
			return err
		}
	}

	return nil
}

func (p *Project) createFiles() error {
	files := map[string]string{
		"cmd/" + p.Name + "/main.go": p.mainTemplate(),
		"internal/config/config.go":  p.configTemplate(),
		"internal/config/otel.go":    p.configOtelTemplate(),
		"deploy/docker/Dockerfile":   p.dockerfileTemplate(),
		".gitignore":                 p.gitignoreTemplate(),
		"Makefile":                   p.makefileTemplate(),
	}

	for path, content := range files {
		fullPath := filepath.Join(p.OutputDir, path)
		if err := os.WriteFile(fullPath, []byte(content), filePerm); err != nil {
			return err
		}
	}

	return nil
}

func (p *Project) initGoMod() error {
	cmd := exec.Command("go", "mod", "init", p.Module)
	cmd.Dir = p.OutputDir
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func (p *Project) goModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.OutputDir
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func (p *Project) mainTemplate() string {
	return fmt.Sprintf(`package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"%s/internal/config"
	"github.com/s4bb4t/zapang/pkg/logger"
	"go.uber.org/zap"
)

const (
	serviceName = "%s"
)

var (
	cfgPath = flag.String("cfg", "./config.yaml", "Config file path")
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	flag.Parse()
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %%s", err)
	}

	l := logger.New(ctx, serviceName, cfg.Logger, nil)
	defer l.Sync()

	l.Info("starting service",
		zap.String("service", serviceName),
		logger.Environment(cfg.Logger.Environment),
	)

	if err := start(ctx, cfg, l); err != nil {
		l.Fatal("failed to start application", zap.Error(err))
	}
}

func start(ctx context.Context, cfg config.Config, log *logger.Logger) error {
	// todo: application

	<-ctx.Done()
	log.Info("shutting down gracefully")
	return nil
}
`, p.Module, p.Name)
}

func (p *Project) configTemplate() string {
	return "package config\n\nimport (\n\t\"os\"\n\n\t\"github.com/go-faster/errors\"\n\t\"github.com/s4bb4t/zapang/pkg/logger\"\n\t\"gopkg.in/yaml.v3\"\n)\n\ntype (\n\tConfig struct {\n\t\tOtel       OpenTelemetry `yaml:\"otel\"`\n\t\tLogger     logger.Config `yaml:\"logger\"`\n\t}\n)\n\nfunc Load(configPath string) (Config, error) {\n\tvar cfg Config\n\n\tfile, err := os.Open(configPath)\n\tif err != nil {\n\t\treturn cfg, errors.Wrap(err, \"open config file\")\n\t}\n\tdefer file.Close()\n\n\tif err := yaml.NewDecoder(file).Decode(&cfg); err != nil {\n\t\treturn cfg, errors.Wrap(err, \"parse config file\")\n\t}\n\n\treturn cfg, nil\n}\n"
}

func (p *Project) configOtelTemplate() string {
	return "package config\n\ntype OpenTelemetry struct {\n\tCollectorPath string `yaml:\"collector_path\"`\n\tServiceName   string `yaml:\"service_name\"`\n}\n"
}

func (p *Project) dockerfileTemplate() string {
	return fmt.Sprintf(`FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/%s ./cmd/%s

FROM alpine:3.19

RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/%s /bin/%s

ENTRYPOINT ["/bin/%s"]
`, p.Name, p.Name, p.Name, p.Name, p.Name)
}

func (p *Project) gitignoreTemplate() string {
	return `# Binaries
bin/
*.exe
*.dll
*.so
*.dylib

# Test
*.test
coverage.out
coverage.html

# IDE
.idea/
.vscode/
*.swp
*.swo

# Environment
.env
.env.local

# OS
.DS_Store
Thumbs.db

# Vendor (optional)
# vendor/

# Build
dist/
`
}

func (p *Project) makefileTemplate() string {
	return fmt.Sprintf(`BINARY_NAME := %s
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build run test lint clean docker

build:
	go build -ldflags="-s -w" -o bin/$(BINARY_NAME) ./cmd/%s

run: build
	./bin/$(BINARY_NAME)

test:
	go test -race -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

docker:
	docker build -t $(BINARY_NAME):$(VERSION) -f deploy/docker/Dockerfile .
`, p.Name, p.Name)
}
