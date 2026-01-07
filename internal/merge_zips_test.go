package internal

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// createTestZip is a helper function to create a zip file with specified content.
// This makes it easy to generate test data.
func createTestZip(t *testing.T, zipPath string, files map[string]string) {
	t.Helper()
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create test zip file %s: %v", zipPath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for name, content := range files {
		f, err := zipWriter.Create(name)
		if err != nil {
			t.Fatalf("Failed to create entry %s in zip: %v", name, err)
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			t.Fatalf("Failed to write content to %s in zip: %v", name, err)
		}
	}
}

// verifyZipContent is a helper function to check the contents of a zip file.
// It ensures the merged zip has the correct files and content.
func verifyZipContent(t *testing.T, zipPath string, expectedFiles map[string]string) {
	t.Helper()
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("Failed to open merged zip file %s for verification: %v", zipPath, err)
	}
	defer r.Close()

	if len(r.File) != len(expectedFiles) {
		t.Fatalf("Expected %d files in merged zip, but found %d", len(expectedFiles), len(r.File))
	}

	for _, f := range r.File {
		expectedContent, ok := expectedFiles[f.Name]
		if !ok {
			t.Fatalf("Unexpected file %s found in merged zip", f.Name)
		}

		rc, err := f.Open()
		if err != nil {
			t.Fatalf("Failed to open file %s in merged zip: %v", f.Name, err)
		}
		defer rc.Close()

		var content bytes.Buffer
		_, err = io.Copy(&content, rc)
		if err != nil {
			t.Fatalf("Failed to read content of file %s: %v", f.Name, err)
		}

		if content.String() != expectedContent {
			t.Errorf("Content mismatch for file %s. Got \"%s\", want \"%s\"", f.Name, content.String(), expectedContent)
		}
		// Remove the file from the map to track which files we've checked
		delete(expectedFiles, f.Name)
	}

	if len(expectedFiles) > 0 {
		for name := range expectedFiles {
			t.Errorf("Expected file %s not found in merged zip", name)
		}
	}
}

func TestMergeZips(t *testing.T) {
	// Create a temporary directory for our test files.
	// This ensures our test is self-contained and cleans up after itself.
	tempDir := t.TempDir()

	// Define file paths
	baseZipPath := filepath.Join(tempDir, "base.zip")
	overlayZipPath := filepath.Join(tempDir, "overlay.zip")
	outputZipPath := filepath.Join(tempDir, "merged.zip")

	// Define content for the zip files
	baseFiles := map[string]string{
		"a.txt":      "content from base",
		"common.txt": "this should be overwritten",
	}
	overlayFiles := map[string]string{
		"b.txt":      "content from overlay",
		"common.txt": "final content from overlay",
	}

	// Define the expected final content after merging
	expectedFiles := map[string]string{
		"a.txt":      "content from base",
		"b.txt":      "content from overlay",
		"common.txt": "final content from overlay", // Overlay content should win
	}

	// Create the test zip files
	createTestZip(t, baseZipPath, baseFiles)
	createTestZip(t, overlayZipPath, overlayFiles)

	// Run the function we want to test
	err := MergeZips(baseZipPath, overlayZipPath, outputZipPath)
	if err != nil {
		t.Fatalf("MergeZips failed: %v", err)
	}

	// Verify that the output zip contains the correct, merged content
	verifyZipContent(t, outputZipPath, expectedFiles)
}
