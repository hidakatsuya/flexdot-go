package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildFlexdot builds the flexdot binary in the given workDir and returns its path.
func buildFlexdot(t *testing.T, workDir string) string {
	bin := filepath.Join(workDir, "flexdot")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd")
	cmd.Dir = filepath.Join("..", "..")
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build flexdot: %v\n%s", err, string(out))
	}
	return bin
}

func TestInstallSymlinkBasic(t *testing.T) {
	workDir := t.TempDir()
	bin := buildFlexdot(t, workDir)

	// Prepare dotfiles dir and file
	dotfilesDir := filepath.Join(workDir, "dotfiles")
	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	dotfile := filepath.Join(dotfilesDir, "myfile.txt")
	if err := os.WriteFile(dotfile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare index.yml
	indexYml := filepath.Join(dotfilesDir, "index.yml")
	indexContent := `myfile.txt: .`
	if err := os.WriteFile(indexYml, []byte(indexContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare home dir
	homeDir := filepath.Join(workDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Run flexdot install
	cmd := exec.Command(bin, "install", "-H", homeDir, "index.yml")
	cmd.Dir = dotfilesDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot install failed: %v\n%s", err, string(out))
	}

	// Check symlink
	linkPath := filepath.Join(homeDir, "myfile.txt")
	fi, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("symlink not created: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("not a symlink: %v", linkPath)
	}
	dest, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	expected := filepath.Join(dotfilesDir, "myfile.txt")
	if dest != expected {
		t.Errorf("symlink points to %s, want %s", dest, expected)
	}
}

func TestInstallAlreadyLinked(t *testing.T) {
	workDir := t.TempDir()
	bin := buildFlexdot(t, workDir)

	// Prepare dotfiles dir and file
	dotfilesDir := filepath.Join(workDir, "dotfiles")
	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	dotfile := filepath.Join(dotfilesDir, "myfile.txt")
	if err := os.WriteFile(dotfile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare index.yml
	indexYml := filepath.Join(dotfilesDir, "index.yml")
	indexContent := `myfile.txt: .`
	if err := os.WriteFile(indexYml, []byte(indexContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare home dir
	homeDir := filepath.Join(workDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 1st install
	cmd := exec.Command(bin, "install", "-H", homeDir, "index.yml")
	cmd.Dir = dotfilesDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot install (1st) failed: %v\n%s", err, string(out))
	}
	// 2nd install (should print "already linked:" and not change the symlink)
	cmd2 := exec.Command(bin, "install", "-H", homeDir, "index.yml")
	cmd2.Dir = dotfilesDir
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot install (2nd) failed: %v\n%s", err, string(out2))
	}
	if !strings.Contains(string(out2), "already linked:") {
		t.Errorf("expected output to contain 'already linked:', got: %s", string(out2))
	}
	// Check symlink still points to the same file
	linkPath := filepath.Join(homeDir, "myfile.txt")
	fi, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("symlink missing after 2nd install: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("not a symlink after 2nd install: %v", linkPath)
	}
	dest, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("failed to read symlink after 2nd install: %v", err)
	}
	expected := filepath.Join(dotfilesDir, "myfile.txt")
	if dest != expected {
		t.Errorf("symlink after 2nd install points to %s, want %s", dest, expected)
	}
}

func TestInstallLinkUpdated(t *testing.T) {
	workDir := t.TempDir()
	bin := buildFlexdot(t, workDir)

	// Prepare dotfiles dir and file (v1)
	dotfilesDir := filepath.Join(workDir, "dotfiles")
	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	dotfile := filepath.Join(dotfilesDir, "myfile.txt")
	if err := os.WriteFile(dotfile, []byte("hello v1"), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare index.yml (v1)
	indexYml := filepath.Join(dotfilesDir, "index.yml")
	indexContentV1 := `myfile.txt: .`
	if err := os.WriteFile(indexYml, []byte(indexContentV1), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare home dir
	homeDir := filepath.Join(workDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 1st install (link to v1)
	cmd := exec.Command(bin, "install", "-H", homeDir, "index.yml")
	cmd.Dir = dotfilesDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot install (1st) failed: %v\n%s", err, string(out))
	}
	linkPath := filepath.Join(homeDir, "myfile.txt")
	fi, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("symlink not created: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("not a symlink: %v", linkPath)
	}
	dest, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	expectedV1 := filepath.Join(dotfilesDir, "myfile.txt")
	if dest != expectedV1 {
		t.Errorf("symlink points to %s, want %s", dest, expectedV1)
	}

	// Temporarily change the symlink to point to a different file
	altDotfile := filepath.Join(dotfilesDir, "altfile.txt")
	if err := os.WriteFile(altDotfile, []byte("alt content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(linkPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(altDotfile, linkPath); err != nil {
		t.Fatal(err)
	}

	// 2nd install (should update the link target, output "link updated:")
	cmd2 := exec.Command(bin, "install", "-H", homeDir, "index.yml")
	cmd2.Dir = dotfilesDir
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot install (2nd) failed: %v\n%s", err, string(out2))
	}
	if !strings.Contains(string(out2), "link updated:") {
		t.Errorf("expected output to contain 'link updated:', got: %s", string(out2))
	}
	// Check symlink points to the correct file
	fi2, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("symlink missing after update: %v", err)
	}
	if fi2.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("not a symlink after update: %v", linkPath)
	}
	dest2, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("failed to read symlink after update: %v", err)
	}
	expectedV2 := filepath.Join(dotfilesDir, "myfile.txt")
	if dest2 != expectedV2 {
		t.Errorf("symlink after update points to %s, want %s", dest2, expectedV2)
	}
}

func TestBackupAndClearBackups(t *testing.T) {
	workDir := t.TempDir()
	bin := buildFlexdot(t, workDir)

	// Prepare dotfiles dir and file
	dotfilesDir := filepath.Join(workDir, "dotfiles")
	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	dotfile := filepath.Join(dotfilesDir, "myfile.txt")
	if err := os.WriteFile(dotfile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare index.yml
	indexYml := filepath.Join(dotfilesDir, "index.yml")
	indexContent := `myfile.txt: .`
	if err := os.WriteFile(indexYml, []byte(indexContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare home dir and create a file that will be backed up
	homeDir := filepath.Join(workDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}
	homeFile := filepath.Join(homeDir, "myfile.txt")
	if err := os.WriteFile(homeFile, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 1st install (should backup the existing file and create a symlink)
	cmd := exec.Command(bin, "install", "-H", homeDir, "index.yml")
	cmd.Dir = dotfilesDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot install (backup test) failed: %v\n%s", err, string(out))
	}
	if !strings.Contains(string(out), "link created:") {
		t.Errorf("expected output to contain 'link created:', got: %s", string(out))
	}
	if !strings.Contains(string(out), "(backup)") {
		t.Errorf("expected output to contain '(backup)', got: %s", string(out))
	}

	// Check that backup directory exists and contains the old file
	backupDir := filepath.Join(dotfilesDir, "backup")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("failed to read backup dir: %v", err)
	}
	found := false
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 14 {
			backupFile := filepath.Join(backupDir, entry.Name(), "myfile.txt")
			if _, err := os.Stat(backupFile); err == nil {
				found = true
				break
			}
		}
	}
	if !found {
		t.Errorf("backup file not found in backup dir")
	}

	// Run clear-backups
	cmd2 := exec.Command(bin, "clear-backups")
	cmd2.Dir = dotfilesDir
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot clear-backups failed: %v\n%s", err, string(out2))
	}
	// Check that backup directory is now empty or does not exist
	entries, err = os.ReadDir(backupDir)
	if err == nil && len(entries) > 0 {
		t.Errorf("backup dir is not empty after clear-backups")
	}
}

func TestInstallWithConfigYml(t *testing.T) {
	workDir := t.TempDir()
	bin := buildFlexdot(t, workDir)

	// Prepare dotfiles dir and file
	dotfilesDir := filepath.Join(workDir, "dotfiles")
	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	dotfile := filepath.Join(dotfilesDir, "myfile.txt")
	if err := os.WriteFile(dotfile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare index.yml
	indexYml := filepath.Join(dotfilesDir, "index.yml")
	indexContent := `myfile.txt: .`
	if err := os.WriteFile(indexYml, []byte(indexContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Prepare config.yml
	homeDir := filepath.Join(workDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}
	configYml := filepath.Join(dotfilesDir, "config.yml")
	configContent := `keep_max_count: 5
home_dir: "` + homeDir + `"
index_yml: index.yml
`
	if err := os.WriteFile(configYml, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Run flexdot install (no args, config.yml used)
	cmd := exec.Command(bin, "install")
	cmd.Dir = dotfilesDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot install (with config.yml) failed: %v\n%s", err, string(out))
	}

	// Check symlink
	linkPath := filepath.Join(homeDir, "myfile.txt")
	fi, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("symlink not created: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("not a symlink: %v", linkPath)
	}
	dest, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	expected := filepath.Join(dotfilesDir, "myfile.txt")
	if dest != expected {
		t.Errorf("symlink points to %s, want %s", dest, expected)
	}
}

func TestInstallMissingArgsAndConfig(t *testing.T) {
	workDir := t.TempDir()
	bin := buildFlexdot(t, workDir)

	// Prepare dotfiles dir (no config.yml, no index.yml)
	dotfilesDir := filepath.Join(workDir, "dotfiles")
	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Run flexdot install (should fail fast)
	cmd := exec.Command(bin, "install")
	cmd.Dir = dotfilesDir
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected install to fail due to missing args/config, but it succeeded")
	}
	if want := "must be specified"; !strings.Contains(string(out), want) {
		t.Errorf("expected error message to contain %q, got: %s", want, string(out))
	}
}
