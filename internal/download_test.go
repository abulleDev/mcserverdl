package internal_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

// mockProgress is a helper to verify that the progress callback is invoked.
func mockProgress(t *testing.T) func(current, total int64) {
	var lastCurrent int64
	return func(current, total int64) {
		if current < lastCurrent {
			t.Errorf("progress current value did not increase: old=%d, new=%d", lastCurrent, current)
		}
		lastCurrent = current
		if total != -1 && current > total {
			t.Errorf("current progress %d exceeds total %d", current, total)
		}
	}
}

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
		case "/cancel":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("start"))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

			// Wait until the context is cancelled by the client side (which closes the connection)
			<-r.Context().Done()
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
		err := internal.Download(context.Background(), server.URL+"/found", destPath, mockProgress(t))
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

		err := internal.Download(context.Background(), server.URL+"/found-no-length", destPath, mockProgress(t))
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

		err := internal.Download(context.Background(), server.URL+"/not-found-path", destPath, mockProgress(t))
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

		err := internal.Download(context.Background(), server.URL+"/found", tempDir, mockProgress(t))
		if err == nil {
			t.Fatal("expected an error for invalid destination path, but got nil")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "cancelled.txt")
		ctx, cancel := context.WithCancel(context.Background())

		err := internal.Download(ctx, server.URL+"/cancel", destPath, func(current, total int64) {
			// Cancel the context as soon as we receive progress
			cancel()
		})
		if err == nil {
			t.Fatal("expected error for cancelled download, but got nil")
		}

		// Verify the file was removed
		if _, err := os.Stat(destPath); !os.IsNotExist(err) {
			t.Errorf("file should have been removed after cancellation, but exists: %s", destPath)
		}
	})
}
