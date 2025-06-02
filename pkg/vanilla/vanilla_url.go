package vanilla

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type detailManifest struct {
	Downloads struct {
		Server *struct {
			URL string `json:"url"`
		} `json:"server"`
	} `json:"downloads"`
}

// DownloadURL returns the download URL for the Minecraft vanilla server JAR for a given game version.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.16.5", "15w14a", "1.18-pre2").
//
// Returns:
//   - string: the direct download URL for the server JAR file.
//   - error: an error if the version is not found or if any HTTP or JSON parsing issues occur.
func DownloadURL(gameVersion string) (string, error) {
	// URL of the version manifest containing all Minecraft vanilla versions
	const url = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

	// Fetch the version manifest from the Mojang API
	versionResponse, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version manifest: %w", err)
	}
	defer versionResponse.Body.Close()

	// Check for a successful HTTP response
	if versionResponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d when fetching version manifest", versionResponse.StatusCode)
	}

	// Decode the JSON manifest into versionManifest struct
	var versionData versionManifest
	if err := json.NewDecoder(versionResponse.Body).Decode(&versionData); err != nil {
		return "", fmt.Errorf("failed to parse version manifest: %w", err)
	}

	// Find the detail URL for the requested game version
	var detailURL string
	for _, version := range versionData.Versions {
		if version.ID == gameVersion {
			detailURL = version.URL
			break
		}
	}

	// Return an error if the version is not found
	if detailURL == "" {
		return "", fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Fetch the detail manifest for the specific version
	detailResponse, err := http.Get(detailURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version detail manifest: %w", err)
	}
	defer detailResponse.Body.Close()

	// Check for a successful HTTP response
	if detailResponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d when fetching version detail manifest", detailResponse.StatusCode)
	}

	// Decode the JSON manifest into detailManifest struct
	var detailData detailManifest
	if err := json.NewDecoder(detailResponse.Body).Decode(&detailData); err != nil {
		return "", fmt.Errorf("failed to parse version detail manifest: %w", err)
	}

	// Return an error if the server download is not available
	if detailData.Downloads.Server == nil {
		return "", fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Return the server JAR download URL
	return detailData.Downloads.Server.URL, nil
}
