package extract

// inspiration
// https://gosamples.dev/unzip-file/
// https://go-cookbook.com/snippets/files/archiving

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"ryvst/pkg/colour"
)

func Unzip(source string, dest string, strip bool) error {
	colour.Info(3, "Extracting %s\n", filepath.Base(source))

	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("unzip, error on zip.openreader: %w", err)
	}
	defer func() { _ = reader.Close() }()

	// 2. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		if err := UnzipFile(f, dest, strip); err != nil {
			return fmt.Errorf("unzip, error unzipfile: %w", err)
		}
	}

	return nil
}

func ParseName(aName string) (string, error) {
	// ZIP and TAR names conventionally use forward slashes.
	name := filepath.Clean(filepath.FromSlash(aName))
	if name == "." || !filepath.IsLocal(name) {
		return "", fmt.Errorf("non-local archive path %q", aName)
	}
	return name, nil
}

func UnzipFile(f *zip.File, dest string, strip bool) error {
	localname, err := ParseName(f.Name)
	if err != nil {
		return fmt.Errorf("parse name error: %w", err)
	}

	// This remove the first folder e.g: ohmyzsh-master/ or yourplugin-main/
	dropBase := stripComponent(localname, strip)
	filePath := filepath.Join(dest, dropBase)

	// For debug
	fmt.Println("Extract:", filePath)

	// 4. Check if file paths are not vulnerable to Zip Slip
	if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return fmt.Errorf("found directory, but failing to create: %w", err)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("zip, basedir, fail to create: %w", err)
	}

	// 6. Handle SYMLINKS
	if f.Mode()&os.ModeSymlink != 0 {
		zippedFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("zip, error on f.open for symlink: %w", err)
		}

		// In zip files, the content of a symlink entry IS the target path string
		targetBytes, err := io.ReadAll(zippedFile)
		_ = zippedFile.Close()
		if err != nil {
			return fmt.Errorf("zip, error reading symlink target: %w", err)
		}

		// Create the symlink pointing to the target, remove old symlink if exist
		if err := os.Symlink(string(targetBytes), filePath); err != nil {
			if os.IsExist(err) {
				// Remove the existing file/symlink so we can recreate it
				if err := os.Remove(filePath); err != nil {
					return fmt.Errorf("failed to remove existing symlink: %w", err)
				}
				// Try creating it again
				if err := os.Symlink(string(targetBytes), filePath); err != nil {
					return fmt.Errorf("failed to recreate symlink: %w", err)
				}
			} else {
				return fmt.Errorf("zip, error creating symlink: %w", err)
			}
		}
		return nil
	}

	// 7. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("zip, error on os.openfile: %w", err)
	}
	defer func() { _ = destinationFile.Close() }()

	// 8. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("zip, error on f.open: %w", err)
	}
	defer func() { _ = zippedFile.Close() }()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return fmt.Errorf("zip, error on io.copy: %w", err)
	}

	return nil
}
