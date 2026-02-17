package scaffold

import "fmt"

func (p *Project) configTemplate() string {
	if p.Mode.HasOpenAPI() && p.Mode.HasGRPC() {
		return p.configTemplateBoth()
	}
	if p.Mode.HasOpenAPI() {
		return p.configTemplateOpenAPI()
	}
	return p.configTemplateGRPC()
}

func (p *Project) configTemplateGRPC() string {
	return "package config\n\nimport (\n\t\"os\"\n\t\"strings\"\n\n\t\"git.web3gate.ru/rkt/vaultage\"\n\t\"github.com/go-faster/errors\"\n\t\"github.com/s4bb4t/srvmon\"\n\t\"github.com/s4bb4t/zapang\"\n)\n\nconst stageKey = \"STAGE\"\n\ntype (\n\tConfig struct {\n\t\t// Core\n\t\tLogger zapang.Config `json:\"logger\" yaml:\"logger\" mapstructure:\"logger\" validate:\"required\"`\n\t\tSrvMon srvmon.Config `json:\"srvmon\" yaml:\"srvmon\" mapstructure:\"srvmon\" validate:\"required\"`\n\n\t\t// External gRPC services\n\t\tExampleService GRPCService `json:\"example_service\" yaml:\"example_service\" mapstructure:\"example_service\"`\n\n\t\t// Server\n\t\tServer ServerConfig `json:\"server\" yaml:\"server\" mapstructure:\"server\" validate:\"required\"`\n\t\t\n\t\t// OpenTelemetry\n\t\tOpenTelemetry OpenTelemetry `json:\"open_telemetry\" yaml:\"open_telemetry\" mapstructure:\"open_telemetry\"`\n\t}\n\n\t// ServerConfig holds gRPC server settings\n\tServerConfig struct {\n\t\tGRPCAddress string `json:\"grpc_address\" yaml:\"grpc_address\" mapstructure:\"grpc_address\" validate:\"required\"`\n\t\tMetricsAddr string `json:\"metrics_addr\" yaml:\"metrics_addr\" mapstructure:\"metrics_addr\" validate:\"required\"`\n\t}\n\n\t// GRPCService holds gRPC client configuration\n\tGRPCService struct {\n\t\tAddress string `json:\"address\" yaml:\"address\" mapstructure:\"address\"`\n\t}\n)\n\nfunc Load() (*Config, error) {\n\tstage := strings.TrimSpace(os.Getenv(stageKey))\n\tif stage == \"\" {\n\t\treturn nil, errors.New(\"STAGE is not set, please set dev, prod or local\")\n\t}\n\n\tcfg := new(Config)\n\tif _, err := vaultage.ReadConfig[Config](cfg, stage); err != nil {\n\t\treturn nil, errors.Wrap(err, \"read config\")\n\t}\n\n\treturn cfg, nil\n}\n"
}

func (p *Project) configTemplateOpenAPI() string {
	return "package config\n\nimport (\n\t\"os\"\n\t\"strings\"\n\n\t\"git.web3gate.ru/rkt/vaultage\"\n\t\"github.com/go-faster/errors\"\n\t\"github.com/s4bb4t/srvmon\"\n\t\"github.com/s4bb4t/zapang\"\n)\n\nconst stageKey = \"STAGE\"\n\ntype (\n\tConfig struct {\n\t\t// Core\n\t\tLogger zapang.Config `json:\"logger\" yaml:\"logger\" mapstructure:\"logger\" validate:\"required\"`\n\t\tSrvMon srvmon.Config `json:\"srvmon\" yaml:\"srvmon\" mapstructure:\"srvmon\" validate:\"required\"`\n\n\t\t// Server\n\t\tServer ServerConfig `json:\"server\" yaml:\"server\" mapstructure:\"server\" validate:\"required\"`\n\t\t\n\t\t// OpenTelemetry\n\t\tOpenTelemetry OpenTelemetry `json:\"open_telemetry\" yaml:\"open_telemetry\" mapstructure:\"open_telemetry\"`\n\t}\n\n\t// ServerConfig holds HTTP server settings\n\tServerConfig struct {\n\t\tHTTPAddress string `json:\"http_address\" yaml:\"http_address\" mapstructure:\"http_address\" validate:\"required\"`\n\t\tMetricsAddr string `json:\"metrics_addr\" yaml:\"metrics_addr\" mapstructure:\"metrics_addr\" validate:\"required\"`\n\t}\n)\n\nfunc Load() (*Config, error) {\n\tstage := strings.TrimSpace(os.Getenv(stageKey))\n\tif stage == \"\" {\n\t\treturn nil, errors.New(\"STAGE is not set, please set dev, prod or local\")\n\t}\n\n\tcfg := new(Config)\n\tif _, err := vaultage.ReadConfig[Config](cfg, stage); err != nil {\n\t\treturn nil, errors.Wrap(err, \"read config\")\n\t}\n\n\treturn cfg, nil\n}\n"
}

func (p *Project) configTemplateBoth() string {
	return "package config\n\nimport (\n\t\"os\"\n\t\"strings\"\n\n\t\"git.web3gate.ru/rkt/vaultage\"\n\t\"github.com/go-faster/errors\"\n\t\"github.com/s4bb4t/srvmon\"\n\t\"github.com/s4bb4t/zapang\"\n)\n\nconst stageKey = \"STAGE\"\n\ntype (\n\tConfig struct {\n\t\t// Core\n\t\tLogger zapang.Config `json:\"logger\" yaml:\"logger\" mapstructure:\"logger\" validate:\"required\"`\n\t\tSrvMon srvmon.Config `json:\"srvmon\" yaml:\"srvmon\" mapstructure:\"srvmon\" validate:\"required\"`\n\n\t\t// External gRPC services\n\t\tExampleService GRPCService `json:\"example_service\" yaml:\"example_service\" mapstructure:\"example_service\"`\n\n\t\t// Server\n\t\tServer ServerConfig `json:\"server\" yaml:\"server\" mapstructure:\"server\" validate:\"required\"`\n\t\t\n\t\t// OpenTelemetry\n\t\tOpenTelemetry OpenTelemetry `json:\"open_telemetry\" yaml:\"open_telemetry\" mapstructure:\"open_telemetry\"`\n\t}\n\n\t// ServerConfig holds server settings\n\tServerConfig struct {\n\t\tGRPCAddress string `json:\"grpc_address\" yaml:\"grpc_address\" mapstructure:\"grpc_address\" validate:\"required\"`\n\t\tHTTPAddress string `json:\"http_address\" yaml:\"http_address\" mapstructure:\"http_address\" validate:\"required\"`\n\t\tMetricsAddr string `json:\"metrics_addr\" yaml:\"metrics_addr\" mapstructure:\"metrics_addr\" validate:\"required\"`\n\t}\n\n\t// GRPCService holds gRPC client configuration\n\tGRPCService struct {\n\t\tAddress string `json:\"address\" yaml:\"address\" mapstructure:\"address\"`\n\t}\n)\n\nfunc Load() (*Config, error) {\n\tstage := strings.TrimSpace(os.Getenv(stageKey))\n\tif stage == \"\" {\n\t\treturn nil, errors.New(\"STAGE is not set, please set dev, prod or local\")\n\t}\n\n\tcfg := new(Config)\n\tif _, err := vaultage.ReadConfig[Config](cfg, stage); err != nil {\n\t\treturn nil, errors.Wrap(err, \"read config\")\n\t}\n\n\treturn cfg, nil\n}\n"
}

func (p *Project) configOtelTemplate() string {
	return "package config\n\ntype OpenTelemetry struct {\n\tCollectorPath string `json:\"collector_path\" yaml:\"collector_path\" mapstructure:\"collector_path\" validate:\"required\"`\n\tServiceName   string `json:\"service_name\" yaml:\"service_name\" mapstructure:\"service_name\" validate:\"required\"`\n}\n"
}

func (p *Project) configFileTemplate() string {
	if p.Mode.HasOpenAPI() && p.Mode.HasGRPC() {
		return fmt.Sprintf(`logger:
  level: debug
  environment: local

srvmon:
  version: v1.0.0
  grpc_address: 0.0.0.0:60011
  http_address: 0.0.0.0:8085

server:
  grpc_address: "0.0.0.0:60006"
  http_address: "0.0.0.0:8080"
  metrics_addr: "0.0.0.0:2112"

example_service:
  address: "localhost:60009"

open_telemetry:
	collector_path: "localhost:0000"
	service_name: %s
`, p.Name)
	}

	if p.Mode.HasOpenAPI() {
		return fmt.Sprintf(`logger:
  level: debug
  environment: local

srvmon:
  version: v1.0.0
  grpc_address: 0.0.0.0:60011
  http_address: 0.0.0.0:8085

server:
  http_address: "0.0.0.0:8080"
  metrics_addr: "0.0.0.0:2112"

open_telemetry:
	collector_path: "localhost:0000"
	service_name: %s
`, p.Name)
	}

	return fmt.Sprintf(`logger:
  level: debug
  environment: local

srvmon:
  version: v1.0.0
  grpc_address: 0.0.0.0:60011
  http_address: 0.0.0.0:8085

server:
  grpc_address: "0.0.0.0:60006"
  metrics_addr: "0.0.0.0:2112"

example_service:
  address: "localhost:60009"

open_telemetry:
	collector_path: "localhost:0000"
	service_name: %s
`, p.Name)
}
