package unix

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"ryvst/pkg/colour"

	"github.com/ulikunitz/xz"
)

// ExtractTarXz extracts a .tar.xz file to a destination directory.
func ExtractTarXz(source, target string) error {
	colour.Info(4, "Extracting %s to %s\n", filepath.Base(source), target)

	// 1. Open the .tar.xz file
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 2. Wrap the file in an XZ decompression reader
	xzReader, err := xz.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %w", err)
	}

	// 3. Wrap the XZ reader in a TAR reader
	tarReader := tar.NewReader(xzReader)

	// 4. Iterate through the files in the tar archive
	for {
		header, err := tarReader.Next()

		// End of archive
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar reading error: %w", err)
		}

		// 🛡️ SECURITY: Prevent Path Traversal (Zip Slip)
		// Ensure the file path doesn't try to escape the target directory (e.g., "../../etc/passwd")
		cleanName := filepath.Clean(header.Name)
		if strings.HasPrefix(cleanName, "..") || strings.Contains(cleanName, string(filepath.Separator)+"..") {
			return fmt.Errorf("illegal file path detected: %s", header.Name)
		}

		targetPath := filepath.Join(target, cleanName)

		// Handle different file types
		switch header.Typeflag {
		case tar.TypeDir: // It's a directory
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

		case tar.TypeReg: // It's a regular file
			// Ensure parent directory exists (tar files don't always list parent dirs first)
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			// Copy the file contents from the tar reader to the new file
			if _, err := io.Copy(outFile, tarReader); err != nil {
				_ = outFile.Close()
				return fmt.Errorf("failed to write file contents: %w", err)
			}
			_ = outFile.Close()

		case tar.TypeSymlink: // It's a symbolic link
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}

		default:
			fmt.Printf("Ignoring unsupported file type %c in %s\n", header.Typeflag, header.Name)
		}
	}

	return nil
}
