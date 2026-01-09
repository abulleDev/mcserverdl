package paper

import (
	"fmt"
	"strconv"

	"github.com/abulleDev/mcserverdl/internal"
)

type buildVersionManifest struct {
	Builds []int `json:"builds"`
}

// ServerVersions fetches the list of all available PaperMC build numbers for a given game version.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.16.5", "1.13-pre7").
//
// Returns:
//   - []string: a slice of build numbers for the specified game version.
//   - error: an error if the game version is not supported or if any HTTP or JSON decoding issues occur.
func (p *Provider) ServerVersions(gameVersion string) ([]string, error) {
	p.Log("Fetching Paper server versions (builds) for %s...", gameVersion)

	// Build manifest URL for the specified game version
	url := fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s", gameVersion)

	// Fetch and decode the build manifest
	var buildData buildVersionManifest
	if err := internal.FetchJSON(url, &buildData); err != nil {
		return nil, fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Return the versions by converting them from int to string.
	builds := make([]string, 0, len(buildData.Builds))
	for _, build := range buildData.Builds {
		builds = append(builds, strconv.Itoa(build))
	}

	p.Log("Fetched %d Paper builds for %s", len(builds), gameVersion)
	return builds, nil
}
