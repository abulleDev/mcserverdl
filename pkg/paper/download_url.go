package paper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type detailManifest struct {
	Downloads struct {
		ServerDefault struct {
			URL string `json:"url"`
		} `json:"server:default"`
	} `json:"downloads"`
}

// DownloadURL returns the download URL for the PaperMC server JAR for a given game version and build number.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.16.5", "1.13-pre7").
//   - serverVersion: the PaperMC build number for the specified version.
//
// Returns:
//   - string: the direct download URL for the PaperMC server JAR file if the build exists.
//   - error: an error if the game version or build number is not found, or if any HTTP or JSON decoding issues occur.
func (p *Provider) DownloadURL(gameVersion, serverVersion string) (string, error) {
	// URL to validate the existence of a specific build
	url := fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds/%s", gameVersion, serverVersion)

	// Send HTTP GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch JSON from %s: %w", url, err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusNotFound:
		// Handle cases where the version or build is not found
		var errorValue struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(response.Body).Decode(&errorValue); err != nil {
			return "", fmt.Errorf("failed to decode error JSON from %s: %w", url, err)
		}

		switch errorValue.Error {
		case "version_not_found":
			return "", fmt.Errorf("unsupported game version: %s", gameVersion)
		case "build_not_found":
			return "", fmt.Errorf("build number %s not found for version %s", gameVersion, gameVersion)
		default:
			return "", fmt.Errorf("unexpected status %d when fetching JSON from %s", response.StatusCode, url)
		}
	case http.StatusOK:
		// Handle successful response
		var versionInfo detailManifest
		if err := json.NewDecoder(response.Body).Decode(&versionInfo); err != nil {
			return "", fmt.Errorf("failed to decode JSON from %s: %w", url, err)
		}
		return versionInfo.Downloads.ServerDefault.URL, nil
	default:
		// Handle other unexpected statuses
		return "", fmt.Errorf("unexpected status %d when fetching JSON from %s", response.StatusCode, url)
	}
}
