package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Temporal TemporalConfig `yaml:"temporal"`
}

type TemporalConfig struct {
	HostPort  string `yaml:"hostPort"`
	TaskQueue string `yaml:"taskQueue"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
