package install

import (
	"fmt"
)

type Options struct {
	DotfilesDir        string
	KeepMaxBackupCount any // nil or int
}

func Run(indexFile, homeDir, dotfilesDir string, keepMaxBackupCount int) error {
	if err := Install(indexFile, homeDir, dotfilesDir, keepMaxBackupCount); err != nil {
		return fmt.Errorf("install failed: %w", err)
	}
	return nil
}
