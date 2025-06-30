package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abulleDev/mcserverdl/cmd/mcserverdl/internal"
	"github.com/abulleDev/mcserverdl/pkg/fabric"
	"github.com/abulleDev/mcserverdl/pkg/forge"
	"github.com/abulleDev/mcserverdl/pkg/neoforge"
	"github.com/abulleDev/mcserverdl/pkg/paper"
	"github.com/abulleDev/mcserverdl/pkg/vanilla"
)

func main() {
	// Disable timestamp and other log prefixes
	log.SetFlags(0)

	// Define command-line flags
	serverType := flag.String("type", "", "Server type (vanilla, paper, forge, fabric, neoforge)")
	gameVersion := flag.String("game", "", "Game version (e.g., 1.21.6, 1.13-pre7, 25w14craftmine)")
	loaderVersion := flag.String("loader", "", "Loader/build version (default latest)")
	path := flag.String("path", "./", "Download path for the server jar")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if *serverType == "" || *gameVersion == "" {
		log.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	// Check if the path is not a file
	info, err := os.Stat(*path)
	if err == nil {
		// Check if it is a file
		if !info.IsDir() {
			log.Fatalf("Error: The specified path '%s' is an existing file. Please provide a directory path", *path)
		}
		// If it is a directory, proceed.
	} else {
		// An error occurred. Check if it is a "not exist" error
		if !os.IsNotExist(err) {
			log.Fatalf("Error checking path '%s': %v", *path, err)
		}
	}

	var url string

	// Call the appropriate download URL function based on the server type
	switch *serverType {
	case "vanilla":
		url, err = vanilla.DownloadURL(*gameVersion)
	case "paper":
		var build int
		if *loaderVersion == "" {
			log.Printf("No build number specified, searching for the latest for %s...", *gameVersion)
			builds, err := paper.Builds(*gameVersion, true)
			if err != nil {
				log.Fatalf("Error fetching latest paper build: %v", err)
			}
			build = builds[0]
			log.Printf("Latest paper build number is %d", build)
		} else {
			build, err = strconv.Atoi(*loaderVersion)
			if err != nil {
				log.Fatalf("Error: invalid build number for paper: %v", err)
			}
		}
		url, err = paper.DownloadURL(*gameVersion, build)
	case "forge":
		if *loaderVersion == "" {
			log.Printf("No loader version specified, searching for the latest for %s...", *gameVersion)
			loaders, err := forge.Loaders(*gameVersion, true)
			if err != nil {
				log.Fatalf("Error fetching latest forge loader: %v", err)
			}
			*loaderVersion = loaders[0]
			log.Printf("Latest forge loader version is %s", *loaderVersion)
		}
		url, err = forge.DownloadURL(*gameVersion, *loaderVersion)
	case "fabric":
		if *loaderVersion == "" {
			log.Printf("No loader version specified, searching for the latest...")
			loaders, err := fabric.Loaders(true)
			if err != nil {
				log.Fatalf("Error fetching latest fabric loader: %v", err)
			}
			*loaderVersion = loaders[0]
			log.Printf("Latest fabric loader version is %s", *loaderVersion)
		}
		url, err = fabric.DownloadURL(*gameVersion, *loaderVersion)
	case "neoforge":
		if *loaderVersion == "" {
			log.Printf("No loader version specified, searching for the latest for %s...", *gameVersion)
			loaders, err := neoforge.Loaders(*gameVersion, true)
			if err != nil {
				log.Fatalf("Error fetching latest neoforge loader: %v", err)
			}
			*loaderVersion = loaders[0]
			log.Printf("Latest neoforge loader version is %s", *loaderVersion)
		}
		url, err = neoforge.DownloadURL(*gameVersion, *loaderVersion)
	default:
		log.Fatalf("Error: Unknown server type '%s'", *serverType)
	}

	// Handle errors and print the result
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create target directory
	if err := os.MkdirAll(*path, 0755); err != nil {
		log.Fatalf("Error: %v", err)
	}

	// The final installation/download logic depends on the server type
	switch *serverType {
	case "forge":
		if strings.HasSuffix(url, ".jar") {
			// Case 1: The URL points to a standard installer JAR
			installerPath := filepath.Join(*path, "installer.jar")
			log.Printf("Downloading Forge installer from %s...", url)
			if err := internal.Download(url, installerPath); err != nil {
				log.Fatalf("Error downloading installer: %v", err)
			}
			log.Println("Installer downloaded. Please run the following command in the installation directory to complete the server setup:")
			log.Println("java -jar installer.jar --installServer")
		} else if strings.HasSuffix(url, ".zip") {
			// Case 2: The URL points to a patch file that needs to be applied to a vanilla server
			patchPath := filepath.Join(*path, "patch.zip")
			vanillaPath := filepath.Join(*path, "vanilla.jar")
			finalJarPath := filepath.Join(*path, "server.jar")

			// Download the patch file
			log.Printf("Downloading Forge patch file from %s...", url)
			if err := internal.Download(url, patchPath); err != nil {
				log.Fatalf("Error downloading patch file: %v", err)
			}

			// Download the corresponding vanilla server JAR
			log.Printf("Downloading vanilla server for %s...", *gameVersion)
			vanillaURL, err := vanilla.DownloadURL(*gameVersion)
			if err != nil {
				log.Fatalf("Error getting vanilla download URL: %v", err)
			}
			if err := internal.Download(vanillaURL, vanillaPath); err != nil {
				log.Fatalf("Error downloading vanilla server: %v", err)
			}

			// Patch the server
			log.Println("Patching vanilla server...")

			if err := internal.MergeZips(vanillaPath, patchPath, finalJarPath); err != nil {
				log.Fatalf("Error: %v", err)
			}

			defer os.Remove(patchPath)
			defer os.Remove(vanillaPath)

			log.Printf("Successfully created Forge server to %s", *path)
		} else {
			log.Fatalf("Unexpected URL format")
		}
	case "neoforge":
		// Download and install the server automatically
		installerPath := filepath.Join(*path, "installer.jar")
		log.Printf("Downloading NeoForge installer from %s...", url)
		if err := internal.Download(url, installerPath); err != nil {
			log.Fatalf("Error downloading installer: %v", err)
		}
		log.Println("Installer downloaded. Please run the following command in the installation directory to complete the server setup:")
		log.Println("java -jar installer.jar --installServer")
	default:
		// For all other server types, just download the file
		finalJarPath := filepath.Join(*path, "server.jar")
		log.Printf("Downloading server from %s...", url)
		if err := internal.Download(url, finalJarPath); err != nil {
			log.Fatalf("Error: %v", err)
		}
		log.Printf("Successfully downloaded server to %s", *path)
	}
}
