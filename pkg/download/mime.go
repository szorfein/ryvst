package download

import (
	"fmt"
	"io"
	"net/http"
)

// Return the type of a file e.g: "application/zip"
func GetFileContentType(url string) (string, error) {
	// 1. Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }() // Always remember to close the body!

	// Check for HTTP errors (e.g., 404 Not Found)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad HTTP status: %s", resp.Status)
	}

	// 2. Read ONLY the first 512 bytes
	// This prevents downloading the entire file into memory
	limitedReader := io.LimitReader(resp.Body, 512)
	buffer, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 3. Detect the content type from the byte slice
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
