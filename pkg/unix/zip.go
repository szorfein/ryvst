package unix

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
)

func Unzip(source string, dest string) error {
	fmt.Println("Unzipping", source, "...")

	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// 2. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		//fmt.Println("Start decompress file", total)
		err := UnzipFile(f, dest)
		if err != nil {
			return err
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

func UnzipFile(f *zip.File, dest string) error {
	localname, err := ParseName(f.Name)
	if err != nil {
		fmt.Println("ParseName Err", err)
	}

	// This remove the first folder e.g: ohmyzsh-master/ or yourplugin-main/
	lsplit := strings.Split(localname, "/")
	dropBase, err := filepath.Rel(lsplit[0], localname)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("default base %s - new %s\n", localname, dropBase)

	var filePath string
	if dropBase == "." {
		filePath = filepath.Join(dest, localname) // TODO, should be removed at the end...
	} else {
		filePath = filepath.Join(dest, dropBase)
	}
	//fmt.Println("INTO Archive with", filePath)

	// 4. Check if file paths are not vulnerable to Zip Slip
	if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer func() {
		err := destinationFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		err := zippedFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
