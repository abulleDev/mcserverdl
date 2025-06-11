package fabric

import "github.com/abulleDev/mcserverdl/internal"

type fabricBuildVersionManifest []struct {
	Version string `json:"version"`
}

// Loaders fetches the list of all available Fabric loader versions from the official FabricMC API.
//
// Parameters:
//   - latestFirst: if true, returns the loader versions with higher versions first. If false, returns the loader versions with lower versions first.
//
// Returns:
//   - []string: a slice of Fabric loader versions (e.g., "0.16.14", "0.15.11").
//   - error: an error if any HTTP or JSON decoding issues occur.
func Loaders(latestFirst bool) ([]string, error) {
	// URL of the version manifest containing all Minecraft fabric loader versions
	const url = "https://meta2.fabricmc.net/v2/versions/loader"

	// Fetch and decode JSON the fabric loader manifest
	var loaderData fabricBuildVersionManifest
	if err := internal.FetchJSON(url, &loaderData); err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(loaderData))
	if latestFirst {
		// Build the slice from first to last (higher versions first)
		for _, version := range loaderData {
			versions = append(versions, version.Version)
		}
	} else {
		// Build the slice from last to first (lower versions first)
		for i := len(loaderData) - 1; i >= 0; i-- {
			versions = append(versions, loaderData[i].Version)
		}
	}

	return versions, nil
}
