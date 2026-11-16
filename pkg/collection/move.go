package collection

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ryvst/pkg/colour"
	"ryvst/pkg/extract"
	"ryvst/pkg/unix"
)

func MimeDetect(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}

	defer func() { _ = file.Close() }()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return http.DetectContentType(buffer), nil
}

func ToDest(cachedir string, name string, todest string, strip bool) error {
	fp := filepath.Join(cachedir, name)
	td := todest

	mimefile, err := MimeDetect(fp)
	if err != nil {
		return err
	}

	//fmt.Println("MIMETYPE", mimefile)
	switch mimefile {
	case "application/zip":
		err := extract.Unzip(fp, td, strip)
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
		if strings.HasSuffix(name, ".tar.xz") {
			err := extract.ExtractTarXz(fp, td, strip)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("MimeType %s not yet supported", mimefile)
	}
	return nil
}
