package vanilla

import (
	"github.com/abulleDev/mcserverdl/internal"
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
//   - error: an error if any HTTP or JSON decoding issues occur.
func Versions(latestFirst bool) ([]string, error) {
	// URL of the version manifest containing all Minecraft vanilla versions
	const url = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

	// Fetch and decode the version manifest
	var versionData versionManifest
	if err := internal.FetchJSON(url, &versionData); err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(versionData.Versions))
	if latestFirst {
		// Build the slice from first to last (higher versions first)
		for _, version := range versionData.Versions {
			versions = append(versions, version.ID)
		}
	} else {
		// Build the slice from last to first (lower versions first)
		for i := len(versionData.Versions) - 1; i >= 0; i-- {
			versions = append(versions, versionData.Versions[i].ID)
		}
	}

	return versions, nil
}
