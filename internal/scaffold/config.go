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
	return "package config\n\nimport (\n\t\"os\"\n\t\"strings\"\n\n\t\"git.web3gate.ru/rkt/metrico/pkg/metrico\"\n\t\"git.web3gate.ru/rkt/trace/pkg/tracer\"\n\t\"git.web3gate.ru/rkt/vaultage\"\n\t\"github.com/go-faster/errors\"\n\t\"github.com/s4bb4t/srvmon\"\n\t\"github.com/s4bb4t/zapang\"\n)\n\nconst stageKey = \"STAGE\"\n\ntype (\n\tConfig struct {\n\t\t// Core\n\t\tDebug         bool           `json:\"debug\" yaml:\"debug\" mapstructure:\"debug\"`\n\t\tLogger        zapang.Config  `json:\"logger\" yaml:\"logger\" mapstructure:\"logger\" validate:\"required\"`\n\t\tSrvMon        srvmon.Config  `json:\"srvmon\" yaml:\"srvmon\" mapstructure:\"srvmon\" validate:\"required\"`\n\t\tOpenTelemetry tracer.Config  `json:\"open_telemetry\" yaml:\"open_telemetry\" mapstructure:\"open_telemetry\" validate:\"required\"`\n\t\tMetrics       metrico.Config `json:\"metrics\" yaml:\"metrics\" mapstructure:\"metrics\" validate:\"required\"`\n\n\t\t// External gRPC services\n\t\tExampleService GRPCService `json:\"example_service\" yaml:\"example_service\" mapstructure:\"example_service\"`\n\n\t\t// Server\n\t\tServer ServerConfig `json:\"server\" yaml:\"server\" mapstructure:\"server\" validate:\"required\"`\n\t}\n\n\t// ServerConfig holds gRPC server settings\n\tServerConfig struct {\n\t\tGRPCAddress string `json:\"grpc_address\" yaml:\"grpc_address\" mapstructure:\"grpc_address\" validate:\"required\"`\n\t}\n\n\t// GRPCService holds gRPC client configuration\n\tGRPCService struct {\n\t\tAddress string `json:\"address\" yaml:\"address\" mapstructure:\"address\"`\n\t}\n)\n\nfunc Load() (*Config, error) {\n\tstage := strings.TrimSpace(os.Getenv(stageKey))\n\tif stage == \"\" {\n\t\treturn nil, errors.New(\"STAGE is not set, please set dev, prod or local\")\n\t}\n\n\tcfg := new(Config)\n\tif _, err := vaultage.ReadConfig[Config](cfg, stage); err != nil {\n\t\treturn nil, errors.Wrap(err, \"read config\")\n\t}\n\n\treturn cfg, nil\n}\n"
}

func (p *Project) configTemplateOpenAPI() string {
	return "package config\n\nimport (\n\t\"os\"\n\t\"strings\"\n\n\t\"git.web3gate.ru/rkt/metrico/pkg/metrico\"\n\t\"git.web3gate.ru/rkt/trace/pkg/tracer\"\n\t\"git.web3gate.ru/rkt/vaultage\"\n\t\"github.com/go-faster/errors\"\n\t\"github.com/s4bb4t/srvmon\"\n\t\"github.com/s4bb4t/zapang\"\n)\n\nconst stageKey = \"STAGE\"\n\ntype (\n\tConfig struct {\n\t\t// Core\n\t\tDebug         bool           `json:\"debug\" yaml:\"debug\" mapstructure:\"debug\"`\n\t\tLogger        zapang.Config  `json:\"logger\" yaml:\"logger\" mapstructure:\"logger\" validate:\"required\"`\n\t\tSrvMon        srvmon.Config  `json:\"srvmon\" yaml:\"srvmon\" mapstructure:\"srvmon\" validate:\"required\"`\n\t\tOpenTelemetry tracer.Config  `json:\"open_telemetry\" yaml:\"open_telemetry\" mapstructure:\"open_telemetry\" validate:\"required\"`\n\t\tMetrics       metrico.Config `json:\"metrics\" yaml:\"metrics\" mapstructure:\"metrics\" validate:\"required\"`\n\n\t\t// Server\n\t\tServer ServerConfig `json:\"server\" yaml:\"server\" mapstructure:\"server\" validate:\"required\"`\n\t}\n\n\t// ServerConfig holds HTTP server settings\n\tServerConfig struct {\n\t\tHTTPAddress string `json:\"http_address\" yaml:\"http_address\" mapstructure:\"http_address\" validate:\"required\"`\n\t}\n)\n\nfunc Load() (*Config, error) {\n\tstage := strings.TrimSpace(os.Getenv(stageKey))\n\tif stage == \"\" {\n\t\treturn nil, errors.New(\"STAGE is not set, please set dev, prod or local\")\n\t}\n\n\tcfg := new(Config)\n\tif _, err := vaultage.ReadConfig[Config](cfg, stage); err != nil {\n\t\treturn nil, errors.Wrap(err, \"read config\")\n\t}\n\n\treturn cfg, nil\n}\n"
}

func (p *Project) configTemplateBoth() string {
	return "package config\n\nimport (\n\t\"os\"\n\t\"strings\"\n\n\t\"git.web3gate.ru/rkt/metrico/pkg/metrico\"\n\t\"git.web3gate.ru/rkt/trace/pkg/tracer\"\n\t\"git.web3gate.ru/rkt/vaultage\"\n\t\"github.com/go-faster/errors\"\n\t\"github.com/s4bb4t/srvmon\"\n\t\"github.com/s4bb4t/zapang\"\n)\n\nconst stageKey = \"STAGE\"\n\ntype (\n\tConfig struct {\n\t\t// Core\n\t\tDebug         bool           `json:\"debug\" yaml:\"debug\" mapstructure:\"debug\"`\n\t\tLogger        zapang.Config  `json:\"logger\" yaml:\"logger\" mapstructure:\"logger\" validate:\"required\"`\n\t\tSrvMon        srvmon.Config  `json:\"srvmon\" yaml:\"srvmon\" mapstructure:\"srvmon\" validate:\"required\"`\n\t\tOpenTelemetry tracer.Config  `json:\"open_telemetry\" yaml:\"open_telemetry\" mapstructure:\"open_telemetry\" validate:\"required\"`\n\t\tMetrics       metrico.Config `json:\"metrics\" yaml:\"metrics\" mapstructure:\"metrics\" validate:\"required\"`\n\n\t\t// External gRPC services\n\t\tExampleService GRPCService `json:\"example_service\" yaml:\"example_service\" mapstructure:\"example_service\"`\n\n\t\t// Server\n\t\tServer ServerConfig `json:\"server\" yaml:\"server\" mapstructure:\"server\" validate:\"required\"`\n\t}\n\n\t// ServerConfig holds server settings\n\tServerConfig struct {\n\t\tGRPCAddress string `json:\"grpc_address\" yaml:\"grpc_address\" mapstructure:\"grpc_address\" validate:\"required\"`\n\t\tHTTPAddress string `json:\"http_address\" yaml:\"http_address\" mapstructure:\"http_address\" validate:\"required\"`\n\t}\n\n\t// GRPCService holds gRPC client configuration\n\tGRPCService struct {\n\t\tAddress string `json:\"address\" yaml:\"address\" mapstructure:\"address\"`\n\t}\n)\n\nfunc Load() (*Config, error) {\n\tstage := strings.TrimSpace(os.Getenv(stageKey))\n\tif stage == \"\" {\n\t\treturn nil, errors.New(\"STAGE is not set, please set dev, prod or local\")\n\t}\n\n\tcfg := new(Config)\n\tif _, err := vaultage.ReadConfig[Config](cfg, stage); err != nil {\n\t\treturn nil, errors.Wrap(err, \"read config\")\n\t}\n\n\treturn cfg, nil\n}\n"
}

func (p *Project) configFileTemplate() string {
	if p.Mode.HasOpenAPI() && p.Mode.HasGRPC() {
		return fmt.Sprintf(`debug: true

logger:
  level: debug
  environment: local

srvmon:
  version: v1.0.0
  grpc_address: :60011
  http_address: :8085

open_telemetry:
  collector_endpoint: localhost:0000
  service_name: %s

metrics:
  service_name: %s
  port: 2112
  path: /metrics

server:
  grpc_address: "0.0.0.0:60006"
  http_address: "0.0.0.0:8080"

example_service:
  address: "localhost:60009"
`, p.Name, p.Name)
	}

	if p.Mode.HasOpenAPI() {
		return fmt.Sprintf(`debug: true

logger:
  level: debug
  environment: local

srvmon:
  version: v1.0.0
  grpc_address: :60011
  http_address: :8085

open_telemetry:
  collector_endpoint: localhost:0000
  service_name: %s

metrics:
  service_name: %s
  port: 2112
  path: /metrics

server:
  http_address: "0.0.0.0:8080"
`, p.Name, p.Name)
	}

	return fmt.Sprintf(`debug: true

logger:
  level: debug
  environment: local

srvmon:
  version: v1.0.0
  grpc_address: :60011
  http_address: :8085

open_telemetry:
  collector_endpoint: localhost:0000
  service_name: %s

metrics:
  service_name: %s
  port: 2112
  path: /metrics

server:
  grpc_address: "0.0.0.0:60006"

example_service:
  address: "localhost:60009"
`, p.Name, p.Name)
}
