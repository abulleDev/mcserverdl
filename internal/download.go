package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ProgressReader implements io.Reader to track download progress.
// It invokes a callback function with the number of bytes read.
type ProgressReader struct {
	io.Reader

	Total      int64
	Current    int64
	OnProgress func(current, total int64)
}

// Read implements the io.Reader interface.
// It reads into p, updates the current progress, and invokes the OnProgress callback.
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if n > 0 {
		pr.Current += int64(n)
		if pr.OnProgress != nil {
			pr.OnProgress(pr.Current, pr.Total)
		}
	}

	return n, err
}

// Download downloads a file from a given URL to a specified path with context support.
// If the download fails, any partially created file at the destination path will be removed.
func Download(ctx context.Context, url, path string, onProgress func(current, total int64)) error {
	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// Send an HTTP GET request to the URL.
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check for a successful HTTP response.
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status code: %d", response.StatusCode)
	}

	// Create the destination file.
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	// Copy the response body to file.
	written, err := io.Copy(out, &ProgressReader{
		Reader:     response.Body,
		Total:      response.ContentLength, // Get the total size of the file from the Content-Length header.
		OnProgress: onProgress,
	})
	if err != nil {
		out.Close()
		os.Remove(path)
		return err
	}

	// Explicitly close the file to check for write errors (e.g. disk full).
	if closeErr := out.Close(); closeErr != nil {
		// If closing fails, ensures the file is removed and returns the error.
		os.Remove(path)
		return fmt.Errorf("failed to close file: %w", closeErr)
	}

	// Ensure the final progress is reported only on success.
	if onProgress != nil {
		onProgress(written, written)
	}

	return nil
}
