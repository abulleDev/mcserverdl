package forge

import (
	"fmt"
	"strings"

	"github.com/abulleDev/mcserverdl/internal"
)

// DownloadURL returns the download URL for the Forge server file for a given game version and loader version.
// It determines the correct URL format based on the game version.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.21.6", "1.7.10-pre4", "1.4").
//   - loaderVersion: the Forge loader version string (e.g., "14.23.4.2720").
//
// Returns:
//   - string: the direct download URL for the Forge server installer/archive file if the versions exist.
//   - error: an error if the game version or loader version is not found, or if any HTTP or JSON decoding issues occur.
func DownloadURL(gameVersion string, loaderVersion string) (string, error) {
	// URL of the version manifest containing all Minecraft forge versions
	const url = "https://files.minecraftforge.net/net/minecraftforge/forge/maven-metadata.json"

	// Fetch and decode the forge loader manifest
	var loaderData map[string][]string
	if err := internal.FetchJSON(url, &loaderData); err != nil {
		return "", err
	}

	// Some game versions have different naming conventions in the manifest
	var forgeStyleVersion string
	switch gameVersion {
	case "1.7.10-pre4":
		forgeStyleVersion = "1.7.10_pre4"
	case "1.4":
		forgeStyleVersion = "1.4.0"
	default:
		forgeStyleVersion = gameVersion
	}

	// Raw loader versions from the manifest (e.g., "1.7.10-10.13.3.1401-1710ls")
	rawLoaderVersions, ok := loaderData[forgeStyleVersion]
	if !ok {
		return "", fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Find the matching loader version and construct the URL
	for i := len(rawLoaderVersions) - 1; i >= 0; i-- {
		// Extract the loader version from the raw string (e.g., "1.7.10-10.13.3.1401-1710ls" -> "10.13.3.1401")
		if strings.Split(rawLoaderVersions[i], "-")[1] == loaderVersion {
			// The file extension varies depending on the game version
			switch gameVersion {
			case
				"1.5.1",
				"1.5",
				"1.4.7",
				"1.4.6",
				"1.4.5",
				"1.4.4",
				"1.4.3",
				"1.4.2",
				"1.4.1",
				"1.4.0",
				"1.3.2":
				// Older versions use "universal.zip"
				return fmt.Sprintf("https://maven.minecraftforge.net/net/minecraftforge/forge/%s/forge-%s-universal.zip", rawLoaderVersions[i], rawLoaderVersions[i]), nil
			case
				"1.2.5",
				"1.2.4",
				"1.2.3",
				"1.1":
				// Very old versions use "server.zip"
				return fmt.Sprintf("https://maven.minecraftforge.net/net/minecraftforge/forge/%s/forge-%s-server.zip", rawLoaderVersions[i], rawLoaderVersions[i]), nil
			default:
				// Modern versions use "installer.jar"
				return fmt.Sprintf("https://maven.minecraftforge.net/net/minecraftforge/forge/%s/forge-%s-installer.jar", rawLoaderVersions[i], rawLoaderVersions[i]), nil
			}
		}
	}

	return "", fmt.Errorf("loader version %s not found for version %s", loaderVersion, gameVersion)
}
