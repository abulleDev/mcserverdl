package paper

import (
	"fmt"

	"github.com/abulleDev/mcserverdl/internal"
)

type paperBuildVersionManifest struct {
	Builds []int `json:"builds"`
}

// Builds fetches the list of all available PaperMC build numbers for a given game version.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.16.5", "1.13-pre7").
//   - latestFirst: if true, returns the builds with higher build numbers first. If false, returns the builds with lower build numbers first.
//
// Returns:
//   - []int: a slice of build numbers for the specified game version.
//   - error: an error if the game version is not supported or if any HTTP or JSON decoding issues occur.
func Builds(gameVersion string, latestFirst bool) ([]int, error) {
	// Build manifest URL for the specified game version
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s", gameVersion)

	// Fetch and decode the build manifest
	var buildData paperBuildVersionManifest
	if err := internal.FetchJSON(url, &buildData); err != nil {
		return nil, fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Reverse the slice (lower versions first)
	if !latestFirst {
		builds := make([]int, 0, len(buildData.Builds))
		for i := len(buildData.Builds) - 1; i >= 0; i-- {
			builds = append(builds, buildData.Builds[i])
		}
		return builds, nil
	}

	// Return the versions as-is (higher versions first)
	return buildData.Builds, nil
}
