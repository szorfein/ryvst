package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"ryvst/pkg/colour"
)

func WriteToDisk(resp *http.Response, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	colour.Info(4, "%s | %d\n", filepath.Base(dest), size)
	return nil
}

func Https(url string, dest string) error {
	colour.Info(3, "Downloading %s...\n", url)
	// TODO: http.Get don't have timeout
	// should we prefer resp, err := http.NewRequest(http.MethodGet, url, nil) ?
	// https://www.gofaq.org/en/how-to-read-an-http-response-body-in-go/
	resp, err := http.Get(url)
	if resp != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	if err != nil {
		return nil
	}
	error := WriteToDisk(resp, dest)
	if error != nil {
		return error
	}
	return nil
}
