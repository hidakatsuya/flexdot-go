package backup

import (
	"os"
	"path/filepath"
	"time"
)

const baseDir = "backup"

// BackupFile moves the file to a new backup directory and returns the backup directory path.
func BackupFile(file string) (string, error) {
	backupDir := filepath.Join(baseDir, time.Now().Format("20060102150405"))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}
	base := filepath.Base(file)
	dest := filepath.Join(backupDir, base)
	if err := os.Rename(file, dest); err != nil {
		return "", err
	}
	return backupDir, nil
}

// RemoveBackupDirIfEmpty removes the backup directory if it's empty.
func RemoveBackupDirIfEmpty(backupDir string) {
	entries, err := os.ReadDir(backupDir)
	if err == nil && len(entries) == 0 {
		os.Remove(backupDir)
	}
}

// RemoveOutdatedBackups removes old backup directories if keepMaxCount is set (>0).
func RemoveOutdatedBackups(keepMaxCount int) {
	if keepMaxCount <= 0 {
		return
	}
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return
	}
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 14 {
			dirs = append(dirs, filepath.Join(baseDir, entry.Name()))
		}
	}
	if len(dirs) <= keepMaxCount {
		return
	}
	// Sort descending (newest first)
	for i := 0; i < len(dirs)-1; i++ {
		for j := i + 1; j < len(dirs); j++ {
			if dirs[i] < dirs[j] {
				dirs[i], dirs[j] = dirs[j], dirs[i]
			}
		}
	}
	// Remove oldest
	for _, dir := range dirs[keepMaxCount:] {
		os.RemoveAll(dir)
	}
}

// ClearAll removes all backup directories.
func ClearAll() error {
	return os.RemoveAll(baseDir)
}
