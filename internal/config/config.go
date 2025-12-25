package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Port                int           `yaml:"port"`
	Algorithm           string        `yaml:"algorithm"`
	LogFile             string        `yaml:"log_file"`
	HealthCheckInterval int           `yaml:"health_check_interval"`
	HealthCheckTimeout  int           `yaml:"health_check_timeout"`
	Backends            []BackendConf `yaml:"backends"`
}

type BackendConf struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
}

func New(filename string) (*Config, error) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(dat, &cfg); err != nil {
		return nil, err
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return nil, fmt.Errorf("Invalid port: %d", cfg.Port)
	}

	if cfg.HealthCheckInterval <= 0 {
		cfg.HealthCheckInterval = 10
	}
	if cfg.HealthCheckTimeout <= 0 {
		cfg.HealthCheckTimeout = 3
	}
	if len(cfg.Backends) == 0 {
		return nil, fmt.Errorf("no backends configured")
	}
	return &cfg, nil
}
