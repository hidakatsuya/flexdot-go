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
func Run(indexFile, homeDir string, opts Options) error {
	if err := Install(indexFile, homeDir, opts); err != nil {
		return fmt.Errorf("install failed: %w", err)
	}
	return nil
}
