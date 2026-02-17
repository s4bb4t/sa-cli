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

type ProjectMode uint8

const (
	ModeGRPC    ProjectMode = 1 << iota // 0b01
	ModeOpenAPI                         // 0b10
)

func (m ProjectMode) HasGRPC() bool    { return m&ModeGRPC != 0 }
func (m ProjectMode) HasOpenAPI() bool { return m&ModeOpenAPI != 0 }

type Project struct {
	Name      string
	Module    string
	OutputDir string
	Mode      ProjectMode
}

func New(name, module, outputDir string, mode ProjectMode) *Project {
	if module == "" {
		module = name
	}
	return &Project{
		Name:      name,
		Module:    module,
		OutputDir: outputDir,
		Mode:      mode,
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
		//{"Installing dependencies", p.goModTidy},
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

		// Internal packages (private)
		"internal/config",
		"internal/infrastructure/database",
		"internal/presentation",

		// Public packages
		"pkg",

		// Deployments
		"deploy/helm/templates",
		"deploy/docker-compose",

		// Scripts
		"scripts",

		// Tests
		"test/integration",
		"test/e2e",
	}

	if p.Mode.HasOpenAPI() {
		dirs = append(dirs,
			"internal/presentation/rest/v1",
			"tools/bundlespec",
		)
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
		"cmd/" + p.Name + "/main.go":                             p.mainTemplate(),
		"cmd/" + p.Name + "/server.go":                           p.serverTemplate(),
		"cmd/" + p.Name + "/services.go":                         p.servicesTemplate(),
		"internal/config/config.go":                              p.configTemplate(),
		"internal/config/otel.go":                                p.configOtelTemplate(),
		"Dockerfile":                                             p.dockerfileTemplate(),
		".gitignore":                                             p.gitignoreTemplate(),
		".gitlab-ci.yml":                                         p.ciTemplate(),
		"Makefile":                                               p.makefileTemplate(),
		"config.yaml":                                            p.configFileTemplate(),
		"deploy/helm/Chart.yaml":                                 p.helmChartTemplate(),
		"deploy/helm/values.yaml":                                p.helmValuesTemplate(),
		"deploy/helm/values-prod.yaml":                           p.helmValuesProdTemplate(),
		"deploy/helm/templates/deployment.yaml":                  p.helmDeploymentTemplate(),
		"deploy/helm/templates/service.yaml":                     p.helmServiceTemplate(),
		"deploy/helm/templates/hpa.yaml":                         p.helmHPATemplate(),
		"deploy/helm/templates/ingress.yaml":                     p.helmIngressTemplate(),
		"deploy/docker-compose/grafana_prom.docker-compose.yaml": p.dockerComposeMonitoringTemplate(),
		"deploy/docker-compose/prometheus.yml":                   p.prometheusConfigTemplate(),
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

//func (p *Project) goModTidy() error {
//	cmd := exec.Command("go", "mod", "tidy")
//	cmd.Dir = p.OutputDir
//	cmd.Stdout = nil
//	cmd.Stderr = nil
//	return cmd.Run()
//}
