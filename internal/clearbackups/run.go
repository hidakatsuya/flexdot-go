package clearbackups

import (
	"fmt"

	"github.com/hidakatsuya/flexdot-go/internal/backup"
)

func Run() error {
	if err := backup.ClearAll(); err != nil {
		return fmt.Errorf("failed to clear backups: %w", err)
	}
	return nil
}
