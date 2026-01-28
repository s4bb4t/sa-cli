package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GRPCService struct {
	Name      string
	Module    string
	OutputDir string
}

func NewGRPC(name, module, outputDir string) *GRPCService {
	if module == "" {
		module = name
	}
	return &GRPCService{
		Name:      name,
		Module:    module,
		OutputDir: outputDir,
	}
}

func (g *GRPCService) Generate() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Creating directories", g.createDirs},
		{"Generating proto contract", g.createProto},
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
		"api/proto/v1",
		"pkg/grpc/" + g.Name + "/v1",
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

func (g *GRPCService) createProto() error {
	protoPath := filepath.Join(g.OutputDir, "api/proto/v1/"+g.Name+".proto")
	return os.WriteFile(protoPath, []byte(g.protoTemplate()), filePerm)
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

	return os.WriteFile(makefilePath, []byte(g.makefileTemplate()), filePerm)
}

func (g *GRPCService) generateProto() error {
	cmd := exec.Command("make", "proto")
	cmd.Dir = g.OutputDir
	return cmd.Run()
}

func (g *GRPCService) protoTemplate() string {
	return fmt.Sprintf(`syntax = "proto3";

package %s.v1;

option go_package = "%s";

service %s {
}
`,
		g.Name,           // package
		g.Module, g.Name, // go_package
	)
}

func (g *GRPCService) makefileTemplate() string {
	return fmt.Sprintf(`BINARY_NAME := %s
SERVICE_NAME := %s

.PHONY: proto build run clean

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/v1/*proto
	mv api/proto/v1/*.pb.go pkg/grpc/$(SERVICE_NAME)/v1

build:
	go build -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

run: build
	./bin/$(BINARY_NAME)

clean:
	rm -rf bin/
`, g.Name, g.Name)
}

func (g *GRPCService) makefileProtoTarget() string {
	return fmt.Sprintf(`
SERVICE_NAME := %s

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/v1/*proto
	mv api/proto/v1/*.pb.go pkg/grpc/$(SERVICE_NAME)/v1
`, g.Name)
}
