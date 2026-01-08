package fabric

import (
	"fmt"
	"net/http"

	"github.com/abulleDev/mcserverdl/internal"
)

type installerVersionManifest []struct {
	Version string `json:"version"`
}

// DownloadURL returns the download URL for the Fabric server JAR for a given game version and loader version.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.21.5", "25w14craftmine", "1.18-pre2").
//   - serverVersion: the Fabric loader version string (e.g., "0.16.14").
//
// Returns:
//   - string: the direct download URL for the Fabric server JAR file if the versions exist.
//   - error: an error if the game version or loader version is not found, or if any HTTP or JSON decoding issues occur.
func (p *Provider) DownloadURL(gameVersion string, serverVersion string) (string, error) {
	// Check Fabric support for the given game version
	checkGameURL := fmt.Sprintf("https://meta2.fabricmc.net/v2/versions/loader/%s", gameVersion)
	response, err := http.Get(checkGameURL)
	if err != nil {
		return "", fmt.Errorf("failed to validate game version: %w", err)
	}
	response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		// Game version is supported
	case http.StatusBadRequest:
		return "", fmt.Errorf("unsupported game version: %s", gameVersion)
	default:
		return "", fmt.Errorf("unexpected status %d while validating game version: %s", response.StatusCode, gameVersion)
	}

	// Check Fabric support for the given server version
	checkServerURL := fmt.Sprintf("https://meta2.fabricmc.net/v2/versions/loader/%s/%s", gameVersion, serverVersion)
	response, err = http.Get(checkServerURL)
	if err != nil {
		return "", fmt.Errorf("failed to validate server version: %w", err)
	}
	response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		// Server version is supported
	case http.StatusBadRequest:
		return "", fmt.Errorf("unsupported server version: %s", serverVersion)
	default:
		return "", fmt.Errorf("unexpected status %d while validating server version: %s", response.StatusCode, serverVersion)
	}

	// Fetch all available installer versions
	const installerURL = "https://meta2.fabricmc.net/v2/versions/installer"
	var installerData installerVersionManifest
	if err := internal.FetchJSON(installerURL, &installerData); err != nil {
		return "", err
	}

	// Use the latest installer version
	latestInstallerVersion := installerData[0].Version

	// Build and return the download URL
	return fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s/%s/%s/server/jar", gameVersion, serverVersion, latestInstallerVersion), nil
}
