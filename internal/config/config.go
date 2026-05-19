package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Source defines a single log source to tail.
type Source struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Format string `yaml:"format"` // json or text
}

// Output defines where filtered logs are written.
type Output struct {
	Type string `yaml:"type"` // stdout or file
	Path string `yaml:"path"` // used when type is file
}

// Filter defines optional log-level filtering.
type Filter struct {
	Levels []string `yaml:"levels"` // e.g. ["error", "warn"]
}

// Config is the top-level logpipe configuration.
type Config struct {
	Sources []Source `yaml:"sources"`
	Output  Output   `yaml:"output"`
	Filter  Filter   `yaml:"filter"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation failed: %w", err)
	}

	return &cfg, nil
}

// validate performs basic sanity checks on the loaded config.
func (c *Config) validate() error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source must be defined")
	}
	for i, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("source[%d]: name is required", i)
		}
		if s.Path == "" {
			return fmt.Errorf("source[%d] %q: path is required", i, s.Name)
		}
		if s.Format != "" && s.Format != "json" && s.Format != "text" {
			return fmt.Errorf("source[%d] %q: format must be 'json' or 'text'", i, s.Name)
		}
	}
	if c.Output.Type != "" && c.Output.Type != "stdout" && c.Output.Type != "file" {
		return fmt.Errorf("output: type must be 'stdout' or 'file'")
	}
	if c.Output.Type == "file" && c.Output.Path == "" {
		return fmt.Errorf("output: path is required when type is 'file'")
	}
	return nil
}
