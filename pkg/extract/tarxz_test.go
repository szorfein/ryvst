package extract

import (
	"archive/tar"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ulikunitz/xz"
)

// createTestTarXz creates a temporary .tar.xz file with the given files and directories.
// It returns the path to the temporary file, which the caller is responsible for cleaning up.
func createTestTarXz(t *testing.T, files map[string]string, dirs []string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test-*.tar.xz")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpName := tmpFile.Name()

	xzWriter, err := xz.NewWriter(tmpFile)
	if err != nil {
		t.Fatalf("failed to create xz writer: %v", err)
	}

	tarWriter := tar.NewWriter(xzWriter)

	// Add directories first
	for _, dir := range dirs {
		header := &tar.Header{
			Name:     dir,
			Mode:     0o755,
			Typeflag: tar.TypeDir,
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatalf("failed to write dir header: %v", err)
		}
	}

	// Add files
	for name, content := range files {
		header := &tar.Header{
			Name:     name,
			Mode:     0o644,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatalf("failed to write file header: %v", err)
		}
		if _, err := tarWriter.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write file content: %v", err)
		}
	}

	if err := tarWriter.Close(); err != nil {
		t.Fatalf("failed to close tar writer: %v", err)
	}
	if err := xzWriter.Close(); err != nil {
		t.Fatalf("failed to close xz writer: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpName
}

func TestExtractTarXz(t *testing.T) {
	t.Run("successful extraction with suffix stripping", func(t *testing.T) {
		files := map[string]string{
			"test-main/file1.txt":     "hello world",
			"test-main/sub/file2.txt": "foo bar",
		}
		dirs := []string{"test-main", "test-main/sub"}

		source := createTestTarXz(t, files, dirs)
		defer func() { _ = os.Remove(source) }()

		target := t.TempDir()

		err := ExtractTarXz(source, target)
		if err != nil {
			t.Fatalf("ExtractTarXz failed: %v", err)
		}

		// Verify file1.txt (suffix should be stripped)
		content, err := os.ReadFile(filepath.Join(target, "file1.txt"))
		if err != nil {
			t.Fatalf("failed to read file1.txt: %v", err)
		}
		if string(content) != "hello world" {
			t.Errorf("expected 'hello world', got '%s'", string(content))
		}

		// Verify file2.txt in subdirectory
		content, err = os.ReadFile(filepath.Join(target, "sub", "file2.txt"))
		if err != nil {
			t.Fatalf("failed to read sub/file2.txt: %v", err)
		}
		if string(content) != "foo bar" {
			t.Errorf("expected 'foo bar', got '%s'", string(content))
		}
	})

	t.Run("non-existent source file", func(t *testing.T) {
		err := ExtractTarXz("/non/existent/file.tar.xz", t.TempDir())
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
		if !strings.Contains(err.Error(), "failed to open file") {
			t.Errorf("expected 'failed to open file' error, got: %v", err)
		}
	})

	t.Run("invalid xz file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "invalid-*.tar.xz")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		if _, err := tmpFile.WriteString("not a valid xz file"); err != nil {
			t.Fatalf("failed to write in temp file: %v", err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("failed to close tmpfile file: %v", err)
		}
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		err = ExtractTarXz(tmpFile.Name(), t.TempDir())
		if err == nil {
			t.Error("expected error for invalid xz file, got nil")
		}
		if !strings.Contains(err.Error(), "failed to create xz reader") {
			t.Errorf("expected 'failed to create xz reader' error, got: %v", err)
		}
	})

	t.Run("valid xz but invalid tar", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "invalid-tar-*.tar.xz")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		xzWriter, err := xz.NewWriter(tmpFile)
		if err != nil {
			t.Fatalf("failed to create xz writer: %v", err)
		}
		_, _ = xzWriter.Write([]byte("this is not a valid tar archive"))
		_ = xzWriter.Close()
		_ = tmpFile.Close()
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		err = ExtractTarXz(tmpFile.Name(), t.TempDir())
		if err == nil {
			t.Error("expected error for invalid tar archive, got nil")
		}
		if !strings.Contains(err.Error(), "tar reading error") {
			t.Errorf("expected 'tar reading error', got: %v", err)
		}
	})
}

func TestUntarFile(t *testing.T) {
	t.Run("extract regular file", func(t *testing.T) {
		dest := t.TempDir()
		content := []byte("hello world")

		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)
		hdr := &tar.Header{
			Name:     "test.txt",
			Mode:     0o644,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}
		_ = tw.WriteHeader(hdr)
		_, _ = tw.Write(content)
		_ = tw.Close()

		tr := tar.NewReader(&buf)
		header, err := tr.Next()
		if err != nil {
			t.Fatalf("failed to get next header: %v", err)
		}

		err = UntarFile(header, tr, dest)
		if err != nil {
			t.Fatalf("UntarFile failed: %v", err)
		}

		extractedContent, err := os.ReadFile(filepath.Join(dest, "test.txt"))
		if err != nil {
			t.Fatalf("failed to read extracted file: %v", err)
		}
		if string(extractedContent) != string(content) {
			t.Errorf("expected '%s', got '%s'", string(content), string(extractedContent))
		}
	})

	t.Run("extract directory", func(t *testing.T) {
		dest := t.TempDir()

		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)
		hdr := &tar.Header{
			Name:     "testdir",
			Mode:     0o755,
			Typeflag: tar.TypeDir,
		}
		_ = tw.WriteHeader(hdr)
		_ = tw.Close()

		tr := tar.NewReader(&buf)
		header, err := tr.Next()
		if err != nil {
			t.Fatalf("failed to get next header: %v", err)
		}

		err = UntarFile(header, tr, dest)
		if err != nil {
			t.Fatalf("UntarFile failed: %v", err)
		}

		info, err := os.Stat(filepath.Join(dest, "testdir"))
		if err != nil {
			t.Fatalf("failed to stat extracted directory: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected extracted path to be a directory")
		}
	})

	t.Run("extract symlink", func(t *testing.T) {
		dest := t.TempDir()

		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)
		hdr := &tar.Header{
			Name:     "testlink",
			Linkname: "target.txt",
			Typeflag: tar.TypeSymlink,
		}
		_ = tw.WriteHeader(hdr)
		_ = tw.Close()

		tr := tar.NewReader(&buf)
		header, err := tr.Next()
		if err != nil {
			t.Fatalf("failed to get next header: %v", err)
		}

		err = UntarFile(header, tr, dest)
		if err != nil {
			t.Fatalf("UntarFile failed: %v", err)
		}

		info, err := os.Lstat(filepath.Join(dest, "testlink"))
		if err != nil {
			t.Fatalf("failed to lstat extracted symlink: %v", err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			t.Error("expected extracted path to be a symlink")
		}
	})

	t.Run("path traversal prevention (Zip Slip)", func(t *testing.T) {
		dest := t.TempDir()
		header := &tar.Header{
			Name:     "../../../etc/passwd",
			Mode:     0o644,
			Size:     11,
			Typeflag: tar.TypeReg,
		}
		reader := tar.NewReader(bytes.NewReader(nil))

		err := UntarFile(header, reader, dest)
		if err == nil {
			t.Error("expected error for path traversal, got nil")
		}
		if !strings.Contains(err.Error(), "illegal file path detected") {
			t.Errorf("expected 'illegal file path detected' error, got: %v", err)
		}
	})

	t.Run("path traversal prevention with separator", func(t *testing.T) {
		dest := t.TempDir()
		header := &tar.Header{
			Name:     "safe/../../../etc/passwd",
			Mode:     0o644,
			Size:     11,
			Typeflag: tar.TypeReg,
		}
		reader := tar.NewReader(bytes.NewReader(nil))

		err := UntarFile(header, reader, dest)
		if err == nil {
			t.Error("expected error for path traversal, got nil")
		}
		if !strings.Contains(err.Error(), "illegal file path detected") {
			t.Errorf("expected 'illegal file path detected' error, got: %v", err)
		}
	})

	t.Run("ignore unsupported file type", func(t *testing.T) {
		dest := t.TempDir()
		header := &tar.Header{
			Name:     "test.fifo",
			Mode:     0o644,
			Typeflag: tar.TypeFifo,
		}
		reader := tar.NewReader(bytes.NewReader(nil))

		// Should not return an error, just ignore and print
		err := UntarFile(header, reader, dest)
		if err != nil {
			t.Fatalf("UntarFile should not return error for unsupported type, got: %v", err)
		}
	})
}

func TestStripComponent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "strip -main suffix",
			input:    "repo-main/file.txt",
			expected: "file.txt",
		},
		{
			name:     "strip -master suffix with nested path",
			input:    "repo-master/dir/file.txt",
			expected: filepath.Join("dir", "file.txt"),
		},
		{
			name:     "no suffix to strip",
			input:    "repo/file.txt",
			expected: "repo/file.txt",
		},
		{
			name:     "single component",
			input:    "file.txt",
			expected: "file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripComponent(tt.input)
			if result != tt.expected {
				t.Errorf("stripComponent(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
