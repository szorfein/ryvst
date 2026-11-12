package collection

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ryvst/pkg/colour"
	"ryvst/pkg/unix"
)

func MimeDetect(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	mimeType := http.DetectContentType(buffer)
	return mimeType
}

func ToDest(collection string, cacheDir string, tasks []Tasks) error {
	for _, task := range tasks {
		if task.Name == "" { // Git
			continue
		}
		fp := filepath.Join(cacheDir, task.Name)
		td := filepath.Join(os.Getenv("HOME"), task.ToDest)

		mimefile := MimeDetect(fp)
		//fmt.Println("MIMETYPE", mimefile)
		switch mimefile {
		case "application/zip":
			colour.Info(3, "Extracting %s to %s\n", filepath.Base(fp), td)
			err := unix.Unzip(fp, td)
			if err != nil {
				return err
			}
		case "font/ttf", "text/xml; charset=utf-8":
			colour.Info(3, "Moving %s to %s\n", filepath.Base(fp), td)
			err := unix.CopyFile(fp, td)
			if err != nil {
				return err
			}
		case "application/x-gzip": // tar.gz
			err := unix.ExtractGZIP(fp, td)
			if err != nil {
				return err
			}
		case "application/octet-stream": // octet-stream is returned when the mimetype is no found...
			if strings.HasSuffix(task.Name, ".tar.xz") {
				err := unix.ExtractTarXz(fp, td)
				if err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("MimeType %s not yet supported", mimefile)
		}
	}
	return nil
}
