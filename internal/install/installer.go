package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hidakatsuya/flexdot-go/internal/backup"
	"gopkg.in/yaml.v3"
)

func Install(indexFile, homeDir string, opts Options) error {
	// Load index YAML
	f, err := os.Open(indexFile)
	if err != nil {
		return fmt.Errorf("failed to open index file: %w", err)
	}
	defer f.Close()

	var idxMap map[string]any
	if err := yaml.NewDecoder(f).Decode(&idxMap); err != nil {
		return fmt.Errorf("failed to decode index yaml: %w", err)
	}

	errs := 0
	for _, entry := range flattenIndex(idxMap) {
		if err := installLink(entry.DotfilePath, entry.HomeFilePath, opts.DotfilesDir, homeDir, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			errs++
		}
	}

	if errs > 0 {
		return fmt.Errorf("encountered %d errors during install", errs)
	}
	return nil
}

func installLink(dotfilePath, homeFilePath, dotfilesDir, homeDir string, opts Options) error {
	dotfile := filepath.Join(dotfilesDir, dotfilePath)
	homeFile := filepath.Join(homeDir, homeFilePath, filepath.Base(dotfile))

	// Ensure dotfile (symlink target) is an absolute path
	dotfileAbs, err := filepath.Abs(dotfile)
	if err != nil {
		return err
	}

	status := &Status{HomeFile: homeFile}

	fi, err := os.Lstat(homeFile)
	if err == nil && fi.Mode()&os.ModeSymlink != 0 {
		linkDest, err := os.Readlink(homeFile)
		if err == nil && linkDest == dotfileAbs {
			status.Result = AlreadyLinked
			OutputLog(homeDir, homeFile, status)
			return nil
		}
		// Remove old symlink and relink
		os.Remove(homeFile)
		if err := os.Symlink(dotfileAbs, homeFile); err != nil {
			return err
		}
		status.Result = LinkUpdated
		OutputLog(homeDir, homeFile, status)
		return nil
	}
	if err == nil && fi.Mode().IsRegular() {
		// Backup and replace
		backupDir, berr := backup.BackupFile(homeFile)
		if berr != nil {
			return berr
		}
		status.Backuped = true
		backup.RemoveBackupDirIfEmpty(backupDir)
	}
	// Ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(homeFile), 0755); err != nil {
		return err
	}
	if err := os.Symlink(dotfileAbs, homeFile); err != nil {
		return err
	}
	status.Result = LinkCreated

	// Remove outdated backups if needed
	if opts.KeepMaxBackupCount != nil {
		if max, ok := opts.KeepMaxBackupCount.(int); ok && max > 0 {
			backup.RemoveOutdatedBackups(max)
		}
	}
	OutputLog(homeDir, homeFile, status)
	return nil
}

// FlattenIndex traverses the index map and returns a slice of dotfile/homefile path pairs.
func flattenIndex(idx map[string]any) []struct{ DotfilePath, HomeFilePath string } {
	var result []struct{ DotfilePath, HomeFilePath string }
	for root, descendants := range idx {
		flattenDescendants(descendants, []string{root}, &result)
	}
	return result
}

func flattenDescendants(descendants any, paths []string, result *[]struct{ DotfilePath, HomeFilePath string }) {
	switch v := descendants.(type) {
	case map[string]any:
		for k, val := range v {
			newPaths := append(paths, k)
			flattenDescendants(val, newPaths, result)
		}
	case string:
		*result = append(*result, struct{ DotfilePath, HomeFilePath string }{
			DotfilePath:  strings.Join(paths, "/"),
			HomeFilePath: v,
		})
	}
}
