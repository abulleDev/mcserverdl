package purpur

import (
	"context"
	"fmt"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

type buildVersionManifest struct {
	Builds struct {
		All []string `json:"all"`
	} `json:"builds"`
}

// ServerVersions fetches the list of all available PurpurMC build numbers for a given game version.
// It uses a default background context.
func (p *Provider) ServerVersions(gameVersion string) ([]string, error) {
	return p.ServerVersionsContext(context.Background(), gameVersion)
}

// ServerVersionsContext fetches the list of all available PurpurMC build numbers for a given game version with context support.
//
// Parameters:
//   - ctx: the context to control the request lifetime.
//   - gameVersion: the Minecraft version string (e.g., "1.21.11", "1.14.1").
//
// Returns:
//   - []string: a slice of build numbers for the specified game version.
//   - error: an error if the game version is not supported or if any HTTP or JSON decoding issues occur.
func (p *Provider) ServerVersionsContext(ctx context.Context, gameVersion string) ([]string, error) {
	p.Log("Fetching Purpur server versions (builds) for %s...", gameVersion)

	// Build manifest URL for the specified game version
	url := fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s", gameVersion)

	// Fetch and decode the build manifest
	var buildData buildVersionManifest
	if err := internal.FetchJSON(ctx, url, &buildData); err != nil {
		return nil, fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Reverse the slice (higher versions first)
	builds := make([]string, 0, len(buildData.Builds.All))
	for i := len(buildData.Builds.All) - 1; i >= 0; i-- {
		builds = append(builds, buildData.Builds.All[i])
	}

	p.Log("Fetched %d Purpur builds for %s", len(builds), gameVersion)
	return builds, nil
}
