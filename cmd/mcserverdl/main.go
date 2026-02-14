package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/abulleDev/mcserverdl/v2/pkg/factory"
)

func main() {
	// Initialize a logger that writes to stdout without timestamps/prefixes.
	logger := log.New(os.Stdout, "", 0)

	// Define command-line flags for server configuration.
	serverType := flag.String("type", "", "Server type (vanilla, paper, forge, fabric, neoforge, purpur)")
	gameVersion := flag.String("game", "", "Game version (e.g., 1.21.6, 1.13-pre7, 25w14craftmine)")
	serverVersion := flag.String("server", "", "Loader/build version (default latest)")
	path := flag.String("path", "./", "Download path for the server jar")
	showVersion := flag.Bool("version", false, "Print the current version")

	// Parse the provided command-line flags.
	flag.Parse()

	// If version flag is set, print version and exit.
	if *showVersion {
		logger.Println("mcserverdl v2.0.0")
		return
	}

	// Validate mandatory flags (type and game version).
	if *serverType == "" || *gameVersion == "" {
		logger.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	// Validate the download path.
	// We check if the path exists and ensure it is not a file.
	info, err := os.Stat(*path)
	if err == nil && !info.IsDir() {
		// The path exists but is a file, which is invalid for a directory argument.
		logger.Fatalf("Error: The specified path '%s' is an existing file. Please provide a directory path", *path)
	}
	if err != nil && !os.IsNotExist(err) {
		// Handle system errors (e.g., permission denied) during path checks.
		logger.Fatalf("Error checking path '%s': %v", *path, err)
	}
	// If the path does not exist, it will be created later via MkdirAll.

	// Initialize the appropriate server provider using the factory.
	provider, err := factory.New(*serverType)
	if err != nil {
		logger.Fatalf("Error: %v", err)
	}

	// Set the logger for the provider to allow consistent logging.
	provider.SetLogger(logger)

	// If server version is not provided, automatically fetch the latest version.
	// Note: Vanilla is excluded here as its logic is handled differently (usually 1:1 with game version).
	if *serverVersion == "" && *serverType != "vanilla" {
		logger.Printf("No server version specified, fetching the latest for %s...", *gameVersion)
		serverVersions, err := provider.ServerVersions(*gameVersion)
		if err != nil {
			logger.Fatalf("Error fetching latest %s server version: %v", *serverType, err)
		}
		if len(serverVersions) == 0 {
			logger.Fatalf("Error: No server versions found for game version %s", *gameVersion)
		}
		*serverVersion = serverVersions[0]
		logger.Printf("Latest %s server version is %s", *serverType, *serverVersion)
	}

	// Create the target directory if it doesn't exist.
	if err := os.MkdirAll(*path, 0755); err != nil {
		logger.Fatalf("Error: %v", err)
	}

	// Execute the download process with a progress callback.
	err = provider.Download(*gameVersion, *serverVersion, *path, func(current, total int64) {
		if total > 0 {
			// If total size is known, display percentage progress.
			fmt.Printf("\rDownloading... %.2f%%", float64(current)/float64(total)*100)
		} else {
			// Fallback for when total size is unknown.
			fmt.Print("\rDownloading...")
		}

		// Clear the progress line once the download completes.
		if total == current {
			fmt.Print("\r\033[K")
		}
	})
	if err != nil {
		logger.Fatalf("Error: %v", err)
	}
}
