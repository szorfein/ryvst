package collection

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "ryvst_collection_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	t.Run("Valid YAML file with all fields", func(t *testing.T) {
		filename := "valid_config.yaml"
		fullpath := filepath.Join(tempDir, filename)

		content := `every_time: 1h30m
execute: /usr/bin/test
tasks:
  - name: task1
    url: http://example.com/1
    to_dest: /tmp/1
  - name: task2
    url: http://example.com/2
    to_dest: /tmp/2
`
		err := os.WriteFile(fullpath, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		config := Read(tempDir, filename)

		if config.EveryTime != 1*time.Hour+30*time.Minute {
			t.Errorf("Expected EveryTime to be 1h30m, got %v", config.EveryTime)
		}
		if config.Execute != "/usr/bin/test" {
			t.Errorf("Expected Execute to be '/usr/bin/test', got '%s'", config.Execute)
		}
		if len(config.Tasks) != 2 {
			t.Fatalf("Expected 2 tasks, got %d", len(config.Tasks))
		}
		if config.Tasks[0].Name != "task1" || config.Tasks[0].Url != "http://example.com/1" || config.Tasks[0].ToDest != "/tmp/1" {
			t.Errorf("Unexpected first task: %+v", config.Tasks[0])
		}
		if config.Tasks[1].Name != "task2" || config.Tasks[1].Url != "http://example.com/2" || config.Tasks[1].ToDest != "/tmp/2" {
			t.Errorf("Unexpected second task: %+v", config.Tasks[1])
		}
	})

	t.Run("Empty YAML file returns defaults", func(t *testing.T) {
		filename := "empty_config.yaml"
		fullpath := filepath.Join(tempDir, filename)

		err := os.WriteFile(fullpath, []byte(""), 0o644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		config := Read(tempDir, filename)

		if config.EveryTime != 24*time.Hour {
			t.Errorf("Expected default EveryTime to be 24h, got %v", config.EveryTime)
		}
		if config.Execute != "" {
			t.Errorf("Expected default Execute to be '', got '%s'", config.Execute)
		}
		if len(config.Tasks) != 0 {
			t.Errorf("Expected default Tasks to be empty, got %+v", config.Tasks)
		}
	})

	t.Run("Partial YAML file merges with defaults", func(t *testing.T) {
		filename := "partial_config.yaml"
		fullpath := filepath.Join(tempDir, filename)

		content := `every_time: 2h`
		err := os.WriteFile(fullpath, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		config := Read(tempDir, filename)

		if config.EveryTime != 2*time.Hour {
			t.Errorf("Expected EveryTime to be 2h, got %v", config.EveryTime)
		}
		if config.Execute != "" {
			t.Errorf("Expected Execute to be '', got '%s'", config.Execute)
		}
	})
}
