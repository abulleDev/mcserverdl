package fabric

import (
	"fmt"
	"net/http"

	"github.com/abulleDev/mcserverdl/internal"
)

type loaderVersionManifest []struct {
	Version string `json:"version"`
}

// ServerVersions fetches the list of all available Fabric loader versions from the official FabricMC API.
// It also verifies that the provided game version is supported by Fabric.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.21.5", "25w14craftmine", "1.18-pre2").
//
// Returns:
//   - []string: a slice of Fabric loader versions (e.g., "0.16.14", "0.15.11").
//   - error: an error if the game version is not supported or if any HTTP or JSON decoding issues occur.
func (p *Provider) ServerVersions(gameVersion string) ([]string, error) {
	p.Log("Fetching Fabric server versions (loaders) for %s...", gameVersion)

	// Check Fabric support for the given version
	// This avoids downloading the large JSON body when we only need to check existence
	checkURL := fmt.Sprintf("https://meta2.fabricmc.net/v2/versions/loader/%s", gameVersion)
	response, err := http.Get(checkURL)
	if err != nil {
		return nil, fmt.Errorf("failed to validate game version: %w", err)
	}
	response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		// Game version is supported
	case http.StatusBadRequest:
		return nil, fmt.Errorf("unsupported game version: %s", gameVersion)
	default:
		return nil, fmt.Errorf("unexpected status %d while validating game version: %s", response.StatusCode, gameVersion)
	}

	// URL of the version manifest containing all Minecraft fabric loader versions
	const url = "https://meta2.fabricmc.net/v2/versions/loader"

	// Fetch and decode JSON the fabric loader manifest
	var loaderData loaderVersionManifest
	if err := internal.FetchJSON(url, &loaderData); err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(loaderData))
	// Build the slice from first to last (higher versions first)
	for _, version := range loaderData {
		versions = append(versions, version.Version)
	}

	p.Log("Fetched %d Fabric loader versions", len(versions))

	return versions, nil
}
