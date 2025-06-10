package fabric

import "github.com/abulleDev/mcserverdl/internal"

type fabricVersionManifest []struct {
	Version string `json:"version"`
}

// Versions fetches the list of all Minecraft Fabric-supported game versions from the official FabricMC API.
//
// Parameters:
//   - latestFirst: if true, returns the versions with higher versions first. If false, returns the versions with lower versions first.
//
// Returns:
//   - []string: a slice of Minecraft versions supported by Fabric (e.g., "1.20.5", "1.18-pre2", "20w51a").
//   - error: an error if any HTTP or JSON decoding issues occur.
func Versions(latestFirst bool) ([]string, error) {
	// URL of the version manifest containing all Minecraft fabric versions
	const url = "https://meta2.fabricmc.net/v2/versions/game"

	// Fetch and decode the fabric version manifest
	var versionData fabricVersionManifest
	if err := internal.FetchJSON(url, &versionData); err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(versionData))
	if latestFirst {
		// Build the slice from first to last (higher versions first)
		for _, version := range versionData {
			versions = append(versions, version.Version)
		}
	} else {
		// Build the slice from last to first (lower versions first)
		for i := len(versionData) - 1; i >= 0; i-- {
			versions = append(versions, versionData[i].Version)
		}
	}

	return versions, nil
}
