package scaffold

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type OpenAPIService struct {
	Name      string
	Module    string
	OutputDir string
	InputDir  string
}

func NewOpenAPI(name, module, outputDir, inputDir string) *OpenAPIService {
	return &OpenAPIService{
		Name:      name,
		Module:    module,
		OutputDir: outputDir,
		InputDir:  inputDir,
	}
}

func (o *OpenAPIService) Generate() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Creating OpenAPI directories", o.createDirs},
		{"Copying OpenAPI specs", o.copySpecs},
		{"Generating .ogen.yml", o.createOgenConfig},
		{"Generating openapi go:generate", o.createGenerateFile},
		{"Generating bundlespec tool", o.createBundlespecTool},
		{"Generating swagger.html", o.createSwaggerHTML},
		{"Generating REST server stub", o.createServerStub},
		{"Generating error handler", o.createErrorHandler},
		{"Appending Makefile targets", o.appendMakefileTargets},
		{"Installing dependencies", o.goModTidy},
		{"Bundling OpenAPI specs", o.runBundleSpec},
		{"Running ogen code generation", o.runOgenGenerate},
		{"Tidying modules", o.goModTidy},
	}

	for _, step := range steps {
		fmt.Printf("  →   %s...\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

func (o *OpenAPIService) createDirs() error {
	dirs := []string{
		"internal/presentation/rest/v1",
		"pkg/openapi/v1",
		"tools/bundlespec",
		"cmd/" + o.Name + "/docs",
		"api/openapi",
	}

	for _, dir := range dirs {
		path := filepath.Join(o.OutputDir, dir)
		if err := os.MkdirAll(path, dirPerm); err != nil {
			return err
		}
	}

	return nil
}

func (o *OpenAPIService) copySpecs() error {
	srcDir := filepath.Join(o.InputDir, "openapi")
	dstDir := filepath.Join(o.OutputDir, "api", "openapi")

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dst := filepath.Join(dstDir, rel)

		if d.IsDir() {
			return os.MkdirAll(dst, dirPerm)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dst, data, filePerm)
	})
}

func (o *OpenAPIService) createOgenConfig() error {
	content := `parser:
  infer_types: true
  allow_remote: true
  depth_limit: 1000
  authentication_schemes:
    - 'bearer'

generator:
  convenient_errors: "on"
  ignore_not_implemented: ["empty schema in request body"]
  features:
    enable:
      - 'paths/server'
      - 'client/request/validation'
      - 'server/response/validation'
      - 'ogen/otel'
      - 'ogen/unimplemented'
      - 'debug/example_tests'
`
	return os.WriteFile(filepath.Join(o.OutputDir, ".ogen.yml"), []byte(content), filePerm)
}

func (o *OpenAPIService) createGenerateFile() error {
	content := fmt.Sprintf(`package v1

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target . --config %s --clean %s
`, filepath.Join("..", "..", "..", ".ogen.yml"),
		filepath.Join("..", "..", "..", "cmd", o.Name, "docs", "v1", "openapi.json"))

	return os.WriteFile(filepath.Join(o.OutputDir, "pkg", "openapi", "v1", "generate.go"), []byte(content), filePerm)
}

func (o *OpenAPIService) createBundlespecTool() error {
	content := `package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

// usage: bundlespec <spec_dir> <out_dir>
//
// Discovers all api versions from <spec_dir>/{version}/openapi.yaml
// and bundles each into <out_dir>/{version}/openapi.json.
func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: bundlespec <spec_dir> <out_dir>\n")
		fmt.Fprintf(os.Stderr, "  spec_dir: directory with versioned specs (e.g. api/openapi)\n")
		fmt.Fprintf(os.Stderr, "  out_dir:  output directory (e.g. cmd/docs)\n")
		os.Exit(1)
	}

	specDir := os.Args[1]
	outDir := os.Args[2]

	matches, err := filepath.Glob(filepath.Join(specDir, "*", "openapi.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "glob: %v\n", err)
		os.Exit(1)
	}

	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "no openapi.yaml found in %s/*/\n", specDir)
		os.Exit(1)
	}

	for _, src := range matches {
		version := filepath.Base(filepath.Dir(src))
		dst := filepath.Join(outDir, version, "openapi.json")

		if err := bundle(src, dst); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", version, err)
			os.Exit(1)
		}

		_, _ = fmt.Fprintf(os.Stdout, "%s: %s -> %s\n", version, src, dst)
	}
}

func bundle(src, dst string) error {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromFile(src)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		fmt.Fprintf(os.Stderr, "  warn: validation: %v\n", err)
	}

	doc.InternalizeRefs(loader.Context, openapi3.DefaultRefNameResolver)

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "  bundled %d bytes\n", len(data))
	return nil
}
`
	return os.WriteFile(filepath.Join(o.OutputDir, "tools", "bundlespec", "main.go"), []byte(content), filePerm)
}

func (o *OpenAPIService) createSwaggerHTML() error {
	return os.WriteFile(filepath.Join(o.OutputDir, "cmd", o.Name, "docs", "swagger.html"), []byte(swaggerHTMLTemplate()), filePerm)
}

func (o *OpenAPIService) createServerStub() error {
	content := fmt.Sprintf(`package v1

import (
	api "%s/pkg/openapi/v1"
	"go.uber.org/zap"
)

type Server struct {
	api.UnimplementedHandler
	logger *zap.Logger
}

func New(log *zap.Logger) *Server {
	return &Server{
		logger: log,
	}
}
`, o.Module)
	return os.WriteFile(filepath.Join(o.OutputDir, "internal", "presentation", "rest", "v1", "server.go"), []byte(content), filePerm)
}

func (o *OpenAPIService) createErrorHandler() error {
	content := fmt.Sprintf(`package v1

import (
	"context"
	"net/http"

	api "%s/pkg/openapi/v1"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) NewError(_ context.Context, err error) *api.ErrorStatusCode {
	st, ok := status.FromError(err)
	if ok {
		return s.handleStatusError(st)
	}

	return s.handleApplicationError(err)
}

func (s *Server) handleStatusError(err *status.Status) *api.ErrorStatusCode {
	s.logger.Error("status error",
		zap.Any("message", err.Message()),
		zap.Error(err.Err()),
		zap.Any("details", err.Details()),
	)

	switch err.Code() {
	case codes.InvalidArgument:
		return &api.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   api.Error{Code: http.StatusBadRequest, Message: err.Message()},
		}

	case codes.DeadlineExceeded:
		return &api.ErrorStatusCode{
			StatusCode: http.StatusRequestTimeout,
			Response:   api.Error{Code: http.StatusRequestTimeout, Message: "request timed out"},
		}

	case codes.NotFound:
		return &api.ErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response:   api.Error{Code: http.StatusNotFound, Message: err.Message()},
		}

	case codes.AlreadyExists:
		return &api.ErrorStatusCode{
			StatusCode: http.StatusConflict,
			Response:   api.Error{Code: http.StatusConflict, Message: err.Message()},
		}

	case codes.PermissionDenied, codes.ResourceExhausted:
		return &api.ErrorStatusCode{
			StatusCode: http.StatusForbidden,
			Response:   api.Error{Code: http.StatusForbidden, Message: err.Message()},
		}

	case codes.Unauthenticated:
		return &api.ErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response:   api.Error{Code: http.StatusUnauthorized, Message: "unauthorized"},
		}

	default:
		return &api.ErrorStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response:   api.Error{Code: http.StatusInternalServerError, Message: "internal server error"},
		}
	}
}

func (s *Server) handleApplicationError(err error) *api.ErrorStatusCode {
	s.logger.Warn("error", zap.Error(err))

	switch {
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
		return &api.ErrorStatusCode{
			StatusCode: http.StatusRequestTimeout,
			Response:   api.Error{Code: http.StatusRequestTimeout, Message: "request timed out"},
		}

	default:
		return &api.ErrorStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response:   api.Error{Code: http.StatusInternalServerError, Message: "internal server error"},
		}
	}
}
`, o.Module)
	return os.WriteFile(filepath.Join(o.OutputDir, "internal", "presentation", "rest", "v1", "error.go"), []byte(content), filePerm)
}

func (o *OpenAPIService) appendMakefileTargets() error {
	makefilePath := filepath.Join(o.OutputDir, "Makefile")

	if _, err := os.Stat(makefilePath); err != nil {
		return errors.New("makefile not found in output directory")
	}

	existing, err := os.ReadFile(makefilePath)
	if err != nil {
		return err
	}

	if strings.Contains(string(existing), "bundle-spec:") {
		return nil
	}

	target := fmt.Sprintf(`
bundle-spec:
	go run tools/bundlespec/main.go api/openapi cmd/%s/docs

generate-ogen: bundle-spec
	go generate ./...
`, o.Name)

	f, err := os.OpenFile(makefilePath, os.O_APPEND|os.O_WRONLY, filePerm)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = f.WriteString(target)
	return err
}

func (o *OpenAPIService) runBundleSpec() error {
	cmd := exec.Command("go", "run", "tools/bundlespec/main.go", "api/openapi", "cmd/"+o.Name+"/docs")
	cmd.Dir = o.OutputDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (o *OpenAPIService) runOgenGenerate() error {
	cmd := exec.Command("go", "generate", "./pkg/openapi/...")
	cmd.Dir = o.OutputDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (o *OpenAPIService) goModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = o.OutputDir
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
