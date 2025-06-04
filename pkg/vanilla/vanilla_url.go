package vanilla

import (
	"fmt"

	"github.com/abulleDev/mcserverdl/internal"
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
//   - error: an error if the version is not found or if any HTTP or JSON decoding issues occur.
func DownloadURL(gameVersion string) (string, error) {
	// URL of the version manifest containing all Minecraft vanilla versions
	const url = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

	// Fetch and decode the version manifest
	var versionData versionManifest
	if err := internal.FetchJSON(url, &versionData); err != nil {
		return "", err
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

	// Fetch and decode the version detail manifest
	var detailData detailManifest
	if err := internal.FetchJSON(url, &detailData); err != nil {
		return "", err
	}

	// Return an error if the server download is not available
	if detailData.Downloads.Server == nil {
		return "", fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Return the server JAR download URL
	return detailData.Downloads.Server.URL, nil
}
