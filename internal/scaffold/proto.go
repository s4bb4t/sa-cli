package scaffold

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GRPCService struct {
	Name      string
	OutputDir string
}

func NewGRPC(name, outputDir string) *GRPCService {
	return &GRPCService{
		Name:      name,
		OutputDir: outputDir,
	}
}

func (g *GRPCService) Generate() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Creating directories", g.createDirs},
		{"Generating Makefile", g.createMakefile},
		{"Generating proto", g.generateProto},
		{"Installing dependencies", g.goModTidy},
	}

	for _, step := range steps {
		fmt.Printf("  →   %s...\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

func (g *GRPCService) goModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = g.OutputDir
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func (g *GRPCService) createDirs() error {
	dirs := []string{
		"internal/presentation/grpc/v1",
	}

	for _, dir := range dirs {
		path := filepath.Join(g.OutputDir, dir)
		if err := os.MkdirAll(path, dirPerm); err != nil {
			return err
		}
	}

	return nil
}

func (g *GRPCService) createMakefile() error {
	makefilePath := filepath.Join(g.OutputDir, "Makefile")

	// Check if Makefile exists, append proto target if it does
	if _, err := os.Stat(makefilePath); err == nil {
		existing, err := os.ReadFile(makefilePath)
		if err != nil {
			return err
		}
		if !strings.Contains(string(existing), "proto:") {
			f, err := os.OpenFile(makefilePath, os.O_APPEND|os.O_WRONLY, filePerm)
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()
			_, err = f.WriteString(g.makefileProtoTarget())
			return err
		}
		return nil
	}

	return errors.New("makefile not found in output directory")
}

func (g *GRPCService) generateProto() error {
	cmd := exec.Command("make", "proto")
	cmd.Dir = g.OutputDir
	return cmd.Run()
}

func (g *GRPCService) makefileProtoTarget() string {
	return `
install-protoc:
	@which protoc > /dev/null || (echo "protoc not found, install: https://github.com/protocolbuffers/protobuf/releases" && exit 1)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto: install-protoc
	@rm -rf pkg/grpc/$(BINARY_NAME)
	@SHARED_INCLUDES=$$(find api/proto/shared -maxdepth 1 -type d -name 'v*' | sort | while read dir; do echo "-I $$dir -I $$dir/deps"; done | tr '\n' ' '); \
	for version in $$(find api/proto -maxdepth 1 -type d -name 'v*' | sort); do \
		ver=$$(basename $$version); \
		echo "Generating proto for $$ver..."; \
		mkdir -p pkg/grpc/$(BINARY_NAME)/$$ver; \
		protoc -I api/proto $$SHARED_INCLUDES \
			--go_out=pkg/grpc/$(BINARY_NAME) --go_opt=paths=source_relative \
			--go-grpc_out=pkg/grpc/$(BINARY_NAME) --go-grpc_opt=paths=source_relative \
			api/proto/$$ver/*.proto; \
	done
	@rm -rf pkg/grpc/$(BINARY_NAME)/shared
	@go mod tidy
`
}
