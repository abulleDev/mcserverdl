package vanilla

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type versionManifest struct {
	Versions []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"versions"`
}

// Versions fetches the list of all Minecraft vanilla versions from the official Mojang API version manifest.
//
// Parameters:
//   - latestFirst: if true, returns the versions with higher versions first. If false, returns the versions with lower versions first.
//
// Returns:
//   - []string: a slice of Minecraft versions (e.g., "1.16.5", "15w14a", "1.18-pre2").
//   - error: an error if any HTTP or JSON parsing issues occur.
func Versions(latestFirst bool) ([]string, error) {
	// URL of the version manifest containing all Minecraft vanilla versions
	const url = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

	// Fetch the manifest from the Mojang API
	versionResponse, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch version manifest: %w", err)
	}
	defer versionResponse.Body.Close()

	// Check for a successful HTTP response
	if versionResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d when fetching version manifest", versionResponse.StatusCode)
	}

	// Decode the JSON manifest into versionManifest struct
	var versionData versionManifest
	if err := json.NewDecoder(versionResponse.Body).Decode(&versionData); err != nil {
		return nil, fmt.Errorf("failed to parse version manifest: %w", err)
	}

	// Create a slice with version ID as a value
	versions := make([]string, 0, len(versionData.Versions))
	for _, version := range versionData.Versions {
		versions = append(versions, version.ID)
	}

	if !latestFirst {
		// Reverse the slice order
		for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
			versions[i], versions[j] = versions[j], versions[i]
		}
	}

	return versions, nil
}
