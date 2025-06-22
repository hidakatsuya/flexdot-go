package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	KeepMaxCount *int   `yaml:"keep_max_count"`
	HomeDir      string `yaml:"home_dir"`
	IndexYml     string `yaml:"index_yml"`
}

func DefaultConfig() Config {
	return Config{
		KeepMaxCount: ptrInt(10),
		HomeDir:      "",
		IndexYml:     "",
	}
}

func ptrInt(v int) *int { return &v }

// LoadConfig loads config.yml from the specified directory.
// If the file does not exist, returns nil and no error.
func LoadConfig(dotfilesDir string) (*Config, error) {
	path := dotfilesDir + "/config.yml"
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open config.yml: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config.yml: %w", err)
	}
	return &cfg, nil
}

// GetKeepMaxCount returns the keepMaxCount value or the default (10) if not set.
func (c *Config) GetKeepMaxCount() int {
	if c == nil || c.KeepMaxCount == nil {
		return *DefaultConfig().KeepMaxCount
	}
	return *c.KeepMaxCount
}
