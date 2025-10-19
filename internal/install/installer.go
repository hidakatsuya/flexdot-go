package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hidakatsuya/flexdot-go/internal/backup"
	"gopkg.in/yaml.v3"
)

type StatusResult int

const (
	AlreadyLinked StatusResult = iota
	LinkUpdated
	LinkCreated
)

type Status struct {
	Result   StatusResult
	Backuped bool
}

type Entry struct {
	DotfilePath  string
	HomeFilePath string
}

func Install(indexFile, homeDir, dotfilesDir string, keepMaxBackupCount int) error {
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
	for _, entry := range flattenIndex(idxMap, dotfilesDir) {
		if err := installLink(entry, dotfilesDir, homeDir, keepMaxBackupCount); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			errs++
		}
	}

	if errs > 0 {
		return fmt.Errorf("encountered %d errors during install", errs)
	}
	return nil
}

func installLink(entry Entry, dotfilesDir, homeDir string, keepMaxBackupCount int) error {
	dotfile := filepath.Join(dotfilesDir, entry.DotfilePath)
	homeFile := filepath.Join(homeDir, entry.HomeFilePath, filepath.Base(dotfile))

	dotfileAbs, err := filepath.Abs(dotfile)
	if err != nil {
		return err
	}

	fi, err := os.Lstat(homeFile)
	switch {
	case err == nil && fi.Mode()&os.ModeSymlink != 0:
		return handleSymlink(homeFile, dotfileAbs, homeDir)
	case err == nil && fi.Mode().IsRegular():
		return handleRegularFile(homeFile, dotfileAbs, homeDir, keepMaxBackupCount)
	case os.IsNotExist(err):
		return handleNotExist(homeFile, dotfileAbs, homeDir)
	default:
		return err
	}
}

func handleSymlink(homeFile, dotfileAbs, homeDir string) error {
	status := &Status{}
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

func handleRegularFile(homeFile, dotfileAbs, homeDir string, keepMaxBackupCount int) error {
	status := &Status{}
	backupDir, berr := backup.BackupFile(homeFile)
	if berr != nil {
		return berr
	}
	status.Backuped = true
	backup.RemoveBackupDirIfEmpty(backupDir)

	if err := os.MkdirAll(filepath.Dir(homeFile), 0755); err != nil {
		return err
	}
	if err := os.Symlink(dotfileAbs, homeFile); err != nil {
		return err
	}
	status.Result = LinkCreated

	if keepMaxBackupCount > 0 {
		backup.RemoveOutdatedBackups(keepMaxBackupCount)
	}
	OutputLog(homeDir, homeFile, status)
	return nil
}

func handleNotExist(homeFile, dotfileAbs, homeDir string) error {
	status := &Status{}
	if err := os.MkdirAll(filepath.Dir(homeFile), 0755); err != nil {
		return err
	}
	if err := os.Symlink(dotfileAbs, homeFile); err != nil {
		return err
	}
	status.Result = LinkCreated
	OutputLog(homeDir, homeFile, status)
	return nil
}

// FlattenIndex traverses the index map and returns a slice of dotfile/homefile path pairs.
func flattenIndex(idx map[string]any, dotfilesDir string) []Entry {
	var result []Entry
	for root, descendants := range idx {
		flattenDescendants(descendants, []string{root}, dotfilesDir, &result)
	}
	return result
}

func flattenDescendants(descendants any, paths []string, dotfilesDir string, result *[]Entry) {
	switch v := descendants.(type) {
	case map[string]any:
		for k, val := range v {
			newPaths := append(paths, k)
			flattenDescendants(val, newPaths, dotfilesDir, result)
		}
	case string:
		hasWildcard := false
		wildcardIndex := -1

		for i, p := range paths {
			if strings.Contains(p, "*") {
				hasWildcard = true
				wildcardIndex = i
				break
			}
		}

		if hasWildcard {
			expandWildcard(paths, wildcardIndex, v, dotfilesDir, result)
		} else {
			*result = append(*result, Entry{
				DotfilePath:  strings.Join(paths, "/"),
				HomeFilePath: v,
			})
		}
	}
}

func expandWildcard(paths []string, wildcardIndex int, homeFilePath string, dotfilesDir string, result *[]Entry) {
	patternPath := strings.Join(paths[:wildcardIndex+1], "/")

	fullPattern := filepath.Join(dotfilesDir, patternPath)

	matches, err := filepath.Glob(fullPattern)
	if err != nil || len(matches) == 0 {
		return
	}

	for _, match := range matches {
		relPath, err := filepath.Rel(dotfilesDir, match)
		if err != nil {
			continue
		}

		matchPath := filepath.ToSlash(relPath)

		if wildcardIndex < len(paths)-1 {
			remainingPaths := paths[wildcardIndex+1:]
			matchPath = matchPath + "/" + strings.Join(remainingPaths, "/")
		}

		*result = append(*result, Entry{
			DotfilePath:  matchPath,
			HomeFilePath: homeFilePath,
		})
	}
}
