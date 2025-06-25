package install

import (
	"fmt"
)

// Options for install command.
type Options struct {
	DotfilesDir        string
	KeepMaxBackupCount any // nil or int
}

// Run executes the install subcommand logic.
func Run(indexFile, homeDir, dotfilesDir string, keepMaxBackupCount int) error {
	if err := Install(indexFile, homeDir, dotfilesDir, keepMaxBackupCount); err != nil {
		return fmt.Errorf("install failed: %w", err)
	}
	return nil
}
