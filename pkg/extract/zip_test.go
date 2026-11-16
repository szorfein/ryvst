package extract

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseName tests the path validation logic to prevent directory traversal (Zip Slip).
func TestParseName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:      "valid simple path",
			input:     "file.txt",
			expected:  "file.txt",
			expectErr: false,
		},
		{
			name:      "valid nested path",
			input:     "dir/subdir/file.txt",
			expected:  filepath.Join("dir", "subdir", "file.txt"),
			expectErr: false,
		},
		{
			name:      "directory traversal attempt",
			input:     "../file.txt",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "complex directory traversal",
			input:     "dir/../../file.txt",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "current directory",
			input:     ".",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "absolute path attempt",
			input:     "/etc/passwd",
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseName(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParseName(%q) error = %v, expectErr %v", tt.input, err, tt.expectErr)
			}
			if result != tt.expected {
				t.Errorf("ParseName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestUnzip tests the main extraction logic with a valid zip file.
func TestUnzip(t *testing.T) {
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "extracted")

	// Create a test zip file with a top-level directory
	err := createTestZip(zipPath, map[string]string{
		"topdir-main/file1.txt":        "content1",
		"topdir-main/subdir/file2.txt": "content2",
	})
	if err != nil {
		t.Fatalf("Failed to create test zip: %v", err)
	}

	err = Unzip(zipPath, extractDir, false)
	if err != nil {
		t.Fatalf("Unzip failed: %v", err)
	}

	// Verify extraction
	// Note: UnzipFile intentionally drops the top-level directory ("topdir")
	file1Path := filepath.Join(extractDir, "file1.txt")
	file2Path := filepath.Join(extractDir, "subdir", "file2.txt")

	content1, err := os.ReadFile(file1Path)
	if err != nil {
		t.Fatalf("Failed to read extracted file1: %v", err)
	}
	if string(content1) != "content1" {
		t.Errorf("file1.txt content = %q, want %q", string(content1), "content1")
	}

	content2, err := os.ReadFile(file2Path)
	if err != nil {
		t.Fatalf("Failed to read extracted file2: %v", err)
	}
	if string(content2) != "content2" {
		t.Errorf("file2.txt content = %q, want %q", string(content2), "content2")
	}
}

// TestUnzipZipSlip ensures that malicious zip files attempting directory traversal are rejected.
func TestUnzipZipSlip(t *testing.T) {
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "slip.zip")
	extractDir := filepath.Join(tempDir, "extracted")

	err := createMaliciousZip(zipPath, "../escaped.txt")
	if err != nil {
		t.Fatalf("Failed to create malicious zip: %v", err)
	}

	err = Unzip(zipPath, extractDir, false)
	if err == nil {
		t.Errorf("Expected error for zip slip, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid file path") && !strings.Contains(err.Error(), "non-local archive path") {
		t.Errorf("Expected zip slip error, got: %v", err)
	}

	// Ensure the escaped file was NOT created in the parent directory
	escapedPath := filepath.Join(tempDir, "escaped.txt")
	if _, err := os.Stat(escapedPath); err == nil {
		t.Errorf("Zip slip succeeded: malicious file was created at %q", escapedPath)
	}
}

// TestUnzipFileNotFound ensures proper error handling when the source zip does not exist.
func TestUnzipFileNotFound(t *testing.T) {
	tempDir := t.TempDir()
	err := Unzip(filepath.Join(tempDir, "nonexistent.zip"), filepath.Join(tempDir, "dest"), false)
	if err == nil {
		t.Errorf("Expected error for nonexistent zip file, got nil")
	}
}

// TestUnzipEmpty tests extraction of a valid but empty zip archive.
func TestUnzipEmpty(t *testing.T) {
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "empty.zip")
	extractDir := filepath.Join(tempDir, "extracted")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create empty zip: %v", err)
	}
	w := zip.NewWriter(f)
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close zip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	err = Unzip(zipPath, extractDir, false)
	if err != nil {
		t.Fatalf("Unzip failed on empty zip: %v", err)
	}
}

// --- Helper Functions ---

func createTestZip(path string, files map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	w := zip.NewWriter(f)
	defer func() { _ = w.Close() }()

	for name, content := range files {
		fw, err := w.Create(name)
		if err != nil {
			return err
		}
		if _, err := fw.Write([]byte(content)); err != nil {
			return err
		}
	}
	return nil
}

func createMaliciousZip(path string, maliciousName string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	w := zip.NewWriter(f)

	// Use a custom header to bypass path cleaning in zip.Create
	header := &zip.FileHeader{
		Name:   maliciousName,
		Method: zip.Deflate,
	}
	fw, err := w.CreateHeader(header)
	if err != nil {
		return err
	}
	if _, err := fw.Write([]byte("malicious content")); err != nil {
		return err
	}

	return w.Close()
}
