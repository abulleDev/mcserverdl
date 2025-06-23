package neoforge

import (
	"fmt"
	"slices"
)

// DownloadURL returns the download URL for the NeoForge server file for a given game version and loader version.
// It determines the correct URL format based on the game version.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.21.6", "25w14craftmine", "1.21").
//   - loaderVersion: the NeoForge loader version string (e.g., "21.0.142-beta", "0.25w14craftmine.5-beta").
//
// Returns:
//   - string: the direct download URL for the NeoForge server installer/archive file if the versions exist.
//   - error: an error if the game version or loader version is not found, or if any HTTP or XML decoding issues occur.
func DownloadURL(gameVersion string, loaderVersion string) (string, error) {
	// Fetch all available loader versions for the given game version.
	loaderVersions, err := Loaders(gameVersion, true)
	if err != nil {
		return "", err
	}

	// Check if the requested loader version exists in the list of available loaders.
	if !slices.Contains(loaderVersions, loaderVersion) {
		return "", fmt.Errorf("loader version %s not found for version %s", loaderVersion, gameVersion)
	}

	return fmt.Sprintf("https://maven.neoforged.net/releases/net/neoforged/neoforge/%s/neoforge-%s-installer.jar", loaderVersion, loaderVersion), nil
}
