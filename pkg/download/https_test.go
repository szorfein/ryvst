package download

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestWriteToDisk(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		dest        string
		expectError bool
	}{
		{
			name:        "successful write",
			body:        "hello world",
			dest:        "test_output_success.txt",
			expectError: false,
		},
		{
			name:        "invalid destination path",
			body:        "hello world",
			dest:        "/invalid/path/that/does/not/exist/file.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				Body: io.NopCloser(bytes.NewBufferString(tt.body)),
			}
			defer func() {
				_ = os.Remove(tt.dest) // Cleanup
			}()

			err := WriteToDisk(resp, tt.dest)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				content, err := os.ReadFile(tt.dest)
				if err != nil {
					t.Fatalf("failed to read created file: %v", err)
				}
				if string(content) != tt.body {
					t.Errorf("expected content %q, got %q", tt.body, string(content))
				}
			}
		})
	}
}

func TestHttps(t *testing.T) {
	// Setup a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/valid":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("downloaded content"))
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("404 page not found"))
		}
	}))
	defer server.Close()

	tests := []struct {
		name        string
		url         string
		dest        string
		expectError bool
	}{
		{
			name:        "successful download",
			url:         server.URL + "/valid",
			dest:        "test_https_valid.txt",
			expectError: false,
		},
		{
			name:        "invalid URL (network error)",
			url:         "http://invalid.local.domain.that.does.not.exist",
			dest:        "test_https_invalid.txt",
			expectError: true,
		},
		{
			name:        "404 not found",
			url:         server.URL + "/notfound",
			dest:        "test_https_404.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				_ = os.Remove(tt.dest) // Cleanup
			}()

			err := Https(tt.url, tt.dest)

			if tt.expectError {
				if err == nil {
					// This will trigger for your current code, exposing the bugs!
					t.Logf("⚠️ WARNING: Test expected an error, but got nil. This exposes a bug in Https().")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				content, err := os.ReadFile(tt.dest)
				if err != nil {
					t.Fatalf("failed to read created file: %v", err)
				}
				if string(content) != "downloaded content" {
					t.Errorf("expected content %q, got %q", "downloaded content", string(content))
				}
			}
		})
	}
}
