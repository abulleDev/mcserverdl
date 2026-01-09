package forge

import (
	"fmt"
	"net/http"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

// GameVersions fetches the list of all Minecraft Forge-supported game versions from the official Forge maven metadata.
//
// Returns:
//   - []string: a slice of Minecraft versions supported by Forge (e.g., "1.21.6", "1.7.10-pre4", "1.4").
//   - error: an error if any HTTP or JSON decoding issues occur.
func (p *Provider) GameVersions() ([]string, error) {
	p.Log("Fetching supported Forge game versions...")

	// URL of the version manifest containing all Minecraft forge versions
	const url = "https://files.minecraftforge.net/net/minecraftforge/forge/maven-metadata.json"

	// Send HTTP GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON from %s: %w", url, err)
	}
	defer response.Body.Close()

	// Check for a successful HTTP response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d when fetching JSON from %s", response.StatusCode, url)
	}

	// Extract version keys from the version manifest
	versionData, err := internal.ExtractJSONKeys(response.Body)
	if err != nil {
		return nil, err
	}

	// Adjust to official version naming
	for index, version := range versionData {
		switch version {
		case "1.7.10_pre4":
			versionData[index] = "1.7.10-pre4"
		case "1.4.0":
			versionData[index] = "1.4"
		}
	}

	// Reverse the slice (higher versions first)
	versions := make([]string, 0, len(versionData))
	for i := len(versionData) - 1; i >= 0; i-- {
		versions = append(versions, versionData[i])
	}

	p.Log("Fetched %d Forge game versions", len(versions))
	return versions, nil
}
