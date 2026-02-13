package purpur

import (
	"context"
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

// DownloadURL returns the download URL for the PurpurMC server JAR for a given game version and build number.
// It uses a default background context.
func (p *Provider) DownloadURL(gameVersion, serverVersion string) (string, error) {
	return p.DownloadURLContext(context.Background(), gameVersion, serverVersion)
}

// DownloadURLContext returns the download URL for the PurpurMC server JAR for a given game version and build number with context support.
//
// Parameters:
//   - ctx: the context to control the request lifetime.
//   - gameVersion: the Minecraft version string (e.g., "1.21.11", "1.14.1").
//   - serverVersion: the PurpurMC build number for the specified version.
//
// Returns:
//   - string: the direct download URL for the PurpurMC server JAR file if the build exists.
//   - error: an error if the game version or build number is not found, or if any HTTP or JSON decoding issues occur.
func (p *Provider) DownloadURLContext(ctx context.Context, gameVersion, serverVersion string) (string, error) {
	p.Log("Fetching download URL for Purpur %s build %s...", gameVersion, serverVersion)

	// URL to validate the existence of a specific build
	url := fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s/%s", gameVersion, serverVersion)

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	// Send HTTP GET request
	response, err := http.DefaultClient.Do(req)
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
		case "version not found":
			return "", fmt.Errorf("unsupported game version: %s", gameVersion)
		case "build not found":
			return "", fmt.Errorf("build number %s not found for version %s", gameVersion, gameVersion)
		default:
			return "", fmt.Errorf("unexpected status %d when fetching JSON from %s", response.StatusCode, url)
		}
	case http.StatusOK:
		// Handle successful response
		serverURL := fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s/%s/download", gameVersion, serverVersion)
		p.Log("Fetched Purpur download URL: %s", serverURL)
		return serverURL, nil
	default:
		// Handle other unexpected statuses
		return "", fmt.Errorf("unexpected status %d when fetching JSON from %s", response.StatusCode, url)
	}
}
