package clearbackups

import (
	"fmt"

	"github.com/hidakatsuya/flexdot-go/internal/backup"
)

// Run executes the clear-backups command logic.
// It removes all backup directories and returns an error if any.
func Run() error {
	if err := backup.ClearAll(); err != nil {
		return fmt.Errorf("failed to clear backups: %w", err)
	}
	return nil
}
