package unix

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"ryvst/pkg/colour"
)

func ExtractGZIP(src string, dest string) error {
	srcBase := filepath.Base(src)
	colour.Info(4, "Extracting %s to %s\n", srcBase, dest)

	MkDAllP(dest, 0o700)

	archive, err := os.Open(src)
	if err != nil {
		return nil
	}
	defer func() {
		_ = archive.Close()
	}()

	uncompressedStream, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer func() {
		_ = uncompressedStream.Close()
	}()

	tarReader := tar.NewReader(uncompressedStream)
	for {
		h, err := tarReader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// called strip_component
		lsplit := strings.Split(h.Name, "/")
		strip, err := filepath.Rel(lsplit[0], h.Name)
		if err != nil {
			return err
		}
		var target string
		if strip == "." {
			target = filepath.Join(dest, h.Name) // TODO, should be removed at the end...
		} else {
			target = filepath.Join(dest, strip)
		}

		// ⚠️ SECURITY: Prevent path traversal attacks (e.g., "../../etc/passwd")
		// target := filepath.Join(dest, h.Name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)+string(os.PathSeparator)) && filepath.Clean(target) != filepath.Clean(dest) {
			return loggableError{"illegal file path: " + h.Name}
		}

		switch h.Typeflag {
		case tar.TypeDir: // Directory
			if err := os.MkdirAll(target, os.FileMode(h.Mode)); err != nil {
				return err
			}
		case tar.TypeReg: // Regular file
			// Create parent directories if they don't exist
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(h.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tarReader); err != nil {
				_ = out.Close()
				return err
			}
			_ = out.Close()
		case tar.TypeSymlink: // Symbolic link
			if PathExist(target) { // TODO, not better to remove it and recreate a new?
				continue
			}
			if err := os.Symlink(h.Linkname, target); err != nil {
				return err
			}
			// 👇 HERE IS HOW YOU DEAL WITH PAX HEADERS 👇
		case tar.TypeXGlobalHeader, tar.TypeXHeader:
			// pax_global_header ('g') and pax_header ('x')
			// We just skip them so they aren't written to disk as files.
			log.Printf("Skipping metadata header: %s", h.Name)
			continue
		default:
			return fmt.Errorf("unsupported TAR entry type %d for %q", h.Typeflag, h.Name)
		}
	}
}

// Custom error type for the path traversal check
type loggableError struct{ msg string }

func (e loggableError) Error() string { return e.msg }
