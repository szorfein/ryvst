package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ryvst/pkg/colour"
)

var DefaultHTTPClient = &http.Client{
	Timeout: 600 * time.Second, // Some Network/Site are slow... Tor power?
}

func WriteToDisk(resp *http.Response, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		fmt.Println("First", dest)
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	size, erro := io.Copy(file, resp.Body)
	if erro != nil {
		fmt.Println("Downloading is too long...", dest)
		return fmt.Errorf("failed to write to file: %w", erro)
	}
	colour.Info(4, "%s | %d bytes written\n", filepath.Base(dest), size)
	return nil
}

func Https(url string, dest string) error {
	colour.Info(3, "Downloading %s...\n", url)
	// TODO: http.Get don't have timeout
	// should we prefer resp, err := http.NewRequest(http.MethodGet, url, nil) ?
	// https://www.gofaq.org/en/how-to-read-an-http-response-body-in-go/
	resp, err := DefaultHTTPClient.Get(url)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil {
		return fmt.Errorf("http request failed: %w", err) // FIXED: was returning nil
	}

	// FIXED: Check for non-2xx status codes before writing to disk
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	error := WriteToDisk(resp, dest)
	if error != nil {
		return fmt.Errorf("failed to write to disk: %w", err)
	}
	return nil
}
