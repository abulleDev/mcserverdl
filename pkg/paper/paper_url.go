package paper

import (
	"fmt"
	"slices"
)

// DownloadURL returns the download URL for the PaperMC server JAR for a given game version and build number.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.16.5", "1.13-pre7").
//   - buildNumber: the PaperMC build number for the specified version.
//
// Returns:
//   - string: the direct download URL for the PaperMC server JAR file if the build exists.
//   - error: an error if the game version or build number is not found, or if any HTTP or JSON decoding issues occur.
func DownloadURL(gameVersion string, buildNumber int) (string, error) {
	buildNumbers, err := Builds(gameVersion, true)
	if err != nil {
		return "", err
	}

	// Check if buildNumber exists in buildNumbers
	found := slices.Contains(buildNumbers, buildNumber)
	if !found {
		return "", fmt.Errorf("build number %d not found for version %s", buildNumber, gameVersion)
	}

	return fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d/downloads/paper-%s-%d.jar", gameVersion, buildNumber, gameVersion, buildNumber), nil
}
