package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// ProgressWriter implements io.Writer to track download progress.
// It sends the number of bytes written to a channel.
type ProgressWriter struct {
	ProgressChan chan int
}

// Write implements the io.Writer interface.
// It writes the byte slice p and sends the number of bytes written to the ProgressChan.
func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.ProgressChan <- n
	return n, nil
}

// Download downloads a file from a given URL to a specified path and shows progress.
func Download(url, path string) error {
	// Send an HTTP GET request to the URL.
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check for a successful HTTP response.
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status code: %d", response.StatusCode)
	}

	// Get the total size of the file from the Content-Length header.
	totalSize := response.ContentLength

	// Create the destination file.
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a channel to send progress updates.
	progressChan := make(chan int)

	// Create a WaitGroup to wait for the progress goroutine to finish.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start a goroutine to read from the progress channel and display the progress.
	go func() {
		defer wg.Done()
		if totalSize > 0 {
			// If the total size is known, calculate and display the percentage.
			var downloadedSize int64
			for chunkSize := range progressChan {
				downloadedSize += int64(chunkSize)
				percentage := float64(downloadedSize) / float64(totalSize) * 100
				fmt.Printf("\rDownloading... %.2f%%", percentage)
			}
		} else {
			// If the total size is unknown, show a static message.
			fmt.Print("Downloading...")
			// Consume from the channel to prevent blocking, but don't print anything.
			for range progressChan {
			}
		}
		fmt.Println("\nDownload complete!")
	}()

	// Create a multi-writer that writes to both the file and the progress writer.
	writer := io.MultiWriter(out, &ProgressWriter{ProgressChan: progressChan})

	// Copy the response body to the multi-writer.
	// This will write the data to the file and send progress updates simultaneously.
	_, err = io.Copy(writer, response.Body)
	if err != nil {
		close(progressChan)
		return err
	}

	// Close the progress channel when the download is complete.
	close(progressChan)

	// Wait for the progress goroutine to finish printing its final message.
	wg.Wait()

	return nil
}
