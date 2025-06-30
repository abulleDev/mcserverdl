package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

// MergeZips combines two zip archives. It overlays the files from overlayZipPath
// on top of baseZipPath. If a file exists in both archives, the one from
// overlayZipPath is used in the final outputZipPath.
func MergeZips(baseZipPath, overlayZipPath, outputZipPath string) error {
	// Create the output zip file that will contain the merged content
	outputZipFile, err := os.Create(outputZipPath)
	if err != nil {
		return err
	}
	defer outputZipFile.Close()

	// Create a new zip writer to write to the output file
	zipWriter := zip.NewWriter(outputZipFile)
	defer zipWriter.Close()

	// A map to keep track of file names from the overlay zip
	// This is used to avoid duplicating files that are present in both zips
	overlayFiles := make(map[string]struct{})

	// Open the overlay zip file for reading
	overlayZip, err := zip.OpenReader(overlayZipPath)
	if err != nil {
		return err
	}
	defer overlayZip.Close()

	// Copy all files from the overlay zip to the output zip
	// These files take precedence
	for _, file := range overlayZip.File {
		// Record the file name to track that it has been added
		overlayFiles[file.Name] = struct{}{}

		// Copy the file from the overlay zip to the new zip
		if err := copyFileToZip(zipWriter, file); err != nil {
			return fmt.Errorf("failed to copy file '%s' from overlay zip: %w", file.Name, err)
		}
	}

	// Open the base zip file for reading.
	baseZip, err := zip.OpenReader(baseZipPath)
	if err != nil {
		return err
	}
	defer baseZip.Close()

	// Copy files from the base zip, but only if they weren't in the overlay zip.
	for _, file := range baseZip.File {
		if _, exists := overlayFiles[file.Name]; !exists {
			if err := copyFileToZip(zipWriter, file); err != nil {
				return fmt.Errorf("failed to copy file '%s' from base zip: %w", file.Name, err)
			}
		}
	}

	return nil
}

// copyFileToZip is a helper function that copies a single file from a source
// zip archive to a destination zip writer
func copyFileToZip(writer *zip.Writer, file *zip.File) error {
	// Open the file inside the zip
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create a new file in the destination zip writer,
	// preserving the original file's header information
	header := file.FileHeader
	destFileWriter, err := writer.CreateHeader(&header)
	if err != nil {
		return err
	}

	// Copy the file content
	_, err = io.Copy(destFileWriter, srcFile)
	return err
}
