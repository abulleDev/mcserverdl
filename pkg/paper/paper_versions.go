package paper

import (
	"github.com/abulleDev/mcserverdl/internal"
)

type paperVersionManifest struct {
	Versions []string `json:"versions"`
}

// Versions fetches the list of all Minecraft paper server versions from the official PaperMC API version manifest.
//
// Parameters:
//   - latestFirst: if true, returns the versions with higher versions first. If false, returns the versions with lower versions first.
//
// Returns:
//   - []string: a slice of Minecraft paper server versions (e.g., "1.16.5", "1.13-pre7").
//   - error: an error if any HTTP or JSON decoding issues occur.
func Versions(latestFirst bool) ([]string, error) {
	// URL of the version manifest containing all Minecraft paper server versions
	const url = "https://api.papermc.io/v2/projects/paper"

	// Fetch and decode the paper server version manifest
	var versionData paperVersionManifest
	if err := internal.FetchJSON(url, &versionData); err != nil {
		return nil, err
	}

	// Return the versions as-is (lower versions first)
	if !latestFirst {
		return versionData.Versions, nil
	}

	// Reverse the slice (higher versions first)
	versions := make([]string, 0, len(versionData.Versions))
	for i := len(versionData.Versions) - 1; i >= 0; i-- {
		versions = append(versions, versionData.Versions[i])
	}
	return versions, nil
}
