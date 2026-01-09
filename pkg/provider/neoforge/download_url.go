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
//   - serverVersion: the NeoForge loader version string (e.g., "21.0.142-beta", "0.25w14craftmine.5-beta").
//
// Returns:
//   - string: the direct download URL for the NeoForge server installer/archive file if the versions exist.
//   - error: an error if the game version or loader version is not found, or if any HTTP or XML decoding issues occur.
func (p *Provider) DownloadURL(gameVersion string, serverVersion string) (string, error) {
	p.Log("Fetching download URL for NeoForge %s loader %s...", gameVersion, serverVersion)

	// Fetch all available loader versions for the given game version.
	loaderVersions, err := p.ServerVersions(gameVersion)
	if err != nil {
		return "", err
	}

	// Check if the requested loader version exists in the list of available loaders.
	if !slices.Contains(loaderVersions, serverVersion) {
		return "", fmt.Errorf("loader version %s not found for version %s", serverVersion, gameVersion)
	}

	serverURL := fmt.Sprintf("https://maven.neoforged.net/releases/net/neoforged/neoforge/%s/neoforge-%s-installer.jar", serverVersion, serverVersion)
	p.Log("Fetched NeoForge download URL: %s", serverURL)
	return serverURL, nil
}
