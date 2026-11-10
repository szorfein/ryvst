package unix

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func PathExist(pathname string) bool {
	_, err := os.Stat(pathname)
	return err == nil
}

func Mkdir(pathname string, permission os.FileMode) {
	if err := os.Mkdir(pathname, permission); err != nil {
		if !os.IsExist(err) {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Directory created successfully.")
	}
}

func MkdirP(pathname string, permission os.FileMode) {
	if !PathExist(pathname) {
		fmt.Println("Creating directory", pathname)
		Mkdir(pathname, permission)
	}
}

func OpenFile(pathname string) *os.File {
	file, err := os.Open(pathname)
	if err != nil {
		fmt.Println("err collection.Read", err)
	}
	return file
}

// usage: unix.MkDAllP("/tmp/h/e/l/l/o", 0o700)
func MkDAllP(longpath string, perm os.FileMode) {
	lp := strings.Split(longpath, "/")
	var record string
	for i, dir := range lp {
		if i == 0 {
			continue
		}
		if i >= 1 {
			np := "/" + dir
			record = record + np
			MkdirP(record, perm)
		}
	}
}

func CopyFile(source string, dest string) error {
	base := filepath.Base(source)
	fullDest := filepath.Join(dest, base)

	MkDAllP(dest, 0o700)

	// Open source file
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer func() {
		_ = sourceFile.Close()
	}()

	// Create destination file
	destFile, err := os.Create(fullDest)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer func() {
		_ = destFile.Close()
	}()

	// Copy the content
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	// Sync to ensure data is written to disk
	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination: %w", err)
	}

	return nil
}
