package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Debug   bool
	Verbose bool
}

var cfg *Config

func init() {
	cfg = &Config{}
}

func Get() *Config {
	return cfg
}

func SetDebug(debug bool) {
	cfg.Debug = debug
}

func SetVerbose(verbose bool) {
	cfg.Verbose = verbose
}

func ConfigDir() string {
	if dir := os.Getenv("SA_CLI_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".sa-cli")
}
