package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownload(t *testing.T) {
	// Test data to be served and checked.
	const testContent = "Hello, this is a test file for download."
	contentLength := fmt.Sprintf("%d", len(testContent))

	// Common test server setup that handles different routes for different scenarios.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/found":
			w.Header().Set("Content-Length", contentLength)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, testContent)
		case "/found-no-length":
			// No Content-Length header for this case.
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, testContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	t.Run("success with known content length", func(t *testing.T) {
		// Create a temporary directory for the download. t.TempDir() handles cleanup.
		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "test.txt")

		// Call the function to be tested.
		err := Download(server.URL+"/found", destPath)
		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}

		// Verify the downloaded file's content.
		content, err := os.ReadFile(destPath)
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		if string(content) != testContent {
			t.Errorf("downloaded content mismatch: got %q, want %q", string(content), testContent)
		}
	})

	t.Run("success with unknown content length", func(t *testing.T) {
		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "test_no_length.txt")

		err := Download(server.URL+"/found-no-length", destPath)
		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}

		// Verify the downloaded file's content.
		content, err := os.ReadFile(destPath)
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		if string(content) != testContent {
			t.Errorf("downloaded content mismatch: got %q, want %q", string(content), testContent)
		}
	})

	t.Run("server returns 404 not found", func(t *testing.T) {
		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "not_found.txt")

		err := Download(server.URL+"/not-found-path", destPath)
		if err == nil {
			t.Fatal("expected an error for 404 response, but got nil")
		}

		// Check that the file was not created on failure.
		if _, err := os.Stat(destPath); !os.IsNotExist(err) {
			t.Errorf("file should not have been created for a failed download, but it exists at %s", destPath)
		}
	})

	t.Run("invalid destination path", func(t *testing.T) {
		// Use a directory as a path, which should cause os.Create to fail.
		tempDir := t.TempDir()

		err := Download(server.URL+"/found", tempDir)
		if err == nil {
			t.Fatal("expected an error for invalid destination path, but got nil")
		}
	})
}
