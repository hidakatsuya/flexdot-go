package initcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hidakatsuya/flexdot-go/internal/config"
	"gopkg.in/yaml.v3"
)

// Run executes the init subcommand logic.
// It creates a config.yml file with default values in the current directory.
func Run() error {
	configPath := filepath.Join(".", "config.yml")

	// Check if config.yml already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config.yml already exists in this directory")
	}

	cfg := config.DefaultConfig()
	yml, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := os.WriteFile(configPath, yml, 0644); err != nil {
		return fmt.Errorf("failed to write config.yml: %w", err)
	}

	fmt.Println("Created config.yml with default values.")
	return nil
}
