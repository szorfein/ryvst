package collection

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestGetWorkdir ensures the workdir is correctly derived from the HOME environment variable.
func TestGetWorkdir(t *testing.T) {
	// t.Setenv automatically restores the original value after the test
	t.Setenv("HOME", "/tmp/test_home_12345")

	expected := "/tmp/test_home_12345/.cache/ryvst"
	if got := getWorkdir(); got != expected {
		t.Errorf("getWorkdir() = %v, want %v", got, expected)
	}
}

func TestCreateCacheDir(t *testing.T) {
	// t.TempDir() creates a temporary directory and automatically cleans it up
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	filename := "my_collection.yml"
	expectedDir := filepath.Join(tempDir, ".cache", "ryvst", "my_collection")

	result := CreateCacheDir(filename)

	if result != expectedDir {
		t.Errorf("CreateCacheDir() = %v, want %v", result, expectedDir)
	}

	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected directory %s to be created", expectedDir)
	}
}

func TestGenHashes(t *testing.T) {
	tempDir := t.TempDir()

	file1 := "test1.txt"
	file2 := "test2.txt"

	os.WriteFile(filepath.Join(tempDir, file1), []byte("hello world"), 0o644)
	os.WriteFile(filepath.Join(tempDir, file2), []byte("foo bar"), 0o644)

	// Include empty string to test the `if file != ""` guard
	files := []string{file1, "", file2}

	hashes, err := GenHashes(tempDir, files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(hashes) != 2 {
		t.Errorf("expected 2 hashes, got %d", len(hashes))
	}

	// Verify deterministic SHA-256 hashes
	expectedHash1 := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	expectedHash2 := "fbc1a9f858ea9e177916964bd88c3d37b91a1e84412765e29950777f265c4b75"

	if hashes[file1] != expectedHash1 {
		t.Errorf("hash for %s = %v, want %v", file1, hashes[file1], expectedHash1)
	}
	if hashes[file2] != expectedHash2 {
		t.Errorf("hash for %s = %v, want %v", file2, hashes[file2], expectedHash2)
	}
}

func TestApplyTimeNow(t *testing.T) {
	tempDir := t.TempDir()

	err := ApplyTimeNow(tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lastTxtPath := filepath.Join(tempDir, "last.txt")
	data, err := os.ReadFile(lastTxtPath)
	if err != nil {
		t.Fatalf("failed to read last.txt: %v", err)
	}

	content := string(data)
	parsedTime, err := time.Parse(time.RFC1123, content)
	if err != nil {
		t.Errorf("failed to parse time from last.txt. Content: '%s', Error: %v", content, err)
	}

	// Verify the recorded time is very recent (allowing a slight negative diff due to execution delay)
	timeDiff := time.Since(parsedTime)
	if timeDiff < -1*time.Second || timeDiff > 5*time.Second {
		t.Errorf("expected recorded time to be within the last 5 seconds, got diff: %v", timeDiff)
	}
}
