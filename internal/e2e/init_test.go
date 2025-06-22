package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInitCreatesConfigYml(t *testing.T) {
	workDir := t.TempDir()
	bin := buildFlexdot(t, workDir)

	cmd := exec.Command(bin, "init")
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("flexdot init failed: %v\nOutput: %s", err, string(out))
	}

	configPath := filepath.Join(workDir, "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("config.yml was not created: %v", err)
	}

	content := string(data)
	if !containsYAMLKey(content, "keep_max_count: 10") {
		t.Errorf("config.yml missing keep_max_count: 10, got:\n%s", content)
	}
	if !containsYAMLKey(content, "home_dir: \"\"") && !containsYAMLKey(content, "home_dir: ''") && !containsYAMLKey(content, "home_dir:") {
		t.Errorf("config.yml missing home_dir (empty), got:\n%s", content)
	}
	if !containsYAMLKey(content, "index_yml: \"\"") && !containsYAMLKey(content, "index_yml: ''") && !containsYAMLKey(content, "index_yml:") {
		t.Errorf("config.yml missing index_yml (empty), got:\n%s", content)
	}
}

func TestInitFailsIfConfigExists(t *testing.T) {
	workDir := t.TempDir()
	bin := buildFlexdot(t, workDir)

	// Create an empty config.yml first
	configPath := filepath.Join(workDir, "config.yml")
	if err := os.WriteFile(configPath, []byte("dummy: true\n"), 0644); err != nil {
		t.Fatalf("failed to create dummy config.yml: %v", err)
	}

	cmd := exec.Command(bin, "init")
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("flexdot init should have failed if config.yml exists")
	}
	if !containsSubstring(string(out), "already exists") {
		t.Errorf("expected error about config.yml already existing, got: %s", string(out))
	}
}

// containsYAMLKey checks if the YAML content contains the given key (and value).
func containsYAMLKey(content, key string) bool {
	return containsLine(content, key)
}

// containsSubstring checks if substr is in s.
func containsSubstring(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && (indexOf(s, substr) >= 0))
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	n := len(s)
	m := len(substr)
	for i := 0; i <= n-m; i++ {
		if s[i:i+m] == substr {
			return i
		}
	}
	return -1
}

// containsLine checks if any line in content matches key exactly (ignoring leading/trailing spaces).
func containsLine(content, key string) bool {
	for _, line := range splitLines(content) {
		if trim(line) == key {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
