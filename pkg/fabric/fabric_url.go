package fabric

import (
	"fmt"
	"slices"

	"github.com/abulleDev/mcserverdl/internal"
)

type fabricInstallerVersionManifest []struct {
	Version string `json:"version"`
}

// DownloadURL returns the download URL for the Fabric server JAR for a given game version and loader version.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.21.5", "25w14craftmine", "1.18-pre2").
//   - loaderVersion: the Fabric loader version string (e.g., "0.16.14").
//
// Returns:
//   - string: the direct download URL for the Fabric server JAR file if the versions exist.
//   - error: an error if the game version or loader version is not found, or if any HTTP or JSON decoding issues occur.
func DownloadURL(gameVersion string, loaderVersion string) (string, error) {
	// Fetch all supported game versions
	gameVersions, err := Versions(true)
	if err != nil {
		return "", err
	}

	// Check if gameVersion exists in gameVersions
	gameVersionFound := slices.Contains(gameVersions, gameVersion)
	if !gameVersionFound {
		return "", fmt.Errorf("game version %s not found", gameVersion)
	}

	// Fetch all supported loader versions
	loaderVersions, err := Loaders(true)
	if err != nil {
		return "", err
	}

	// Check if loaderVersion exists in loaderVersions
	loaderVersionFound := slices.Contains(loaderVersions, loaderVersion)
	if !loaderVersionFound {
		return "", fmt.Errorf("loader version %s not found", loaderVersion)
	}

	// Fetch all available installer versions
	const url = "https://meta2.fabricmc.net/v2/versions/installer"
	var installerData fabricInstallerVersionManifest
	if err := internal.FetchJSON(url, &installerData); err != nil {
		return "", err
	}

	// Use the latest installer version
	latestInstallerVersion := installerData[0].Version

	// Build and return the download URL
	return fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s/%s/%s/server/jar", gameVersion, loaderVersion, latestInstallerVersion), nil
}
