package purpur

import (
	"context"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

type versionManifest struct {
	Versions []string `json:"versions"`
}

// GameVersions fetches the list of all Minecraft Purpur-supported game versions from the official PurpurMC API.
// It uses a default background context.
func (p *Provider) GameVersions() ([]string, error) {
	return p.GameVersionsContext(context.Background())
}

// GameVersionsContext fetches the list of all Minecraft Purpur-supported game versions from the official PurpurMC API with context support.
//
// Parameters:
//   - ctx: the context to control the request lifetime.
//
// Returns:
//   - []string: a slice of Minecraft versions supported by Purpur (e.g., "1.21.11", "1.14.1").
//   - error: an error if any HTTP or JSON decoding issues occur.
func (p *Provider) GameVersionsContext(ctx context.Context) ([]string, error) {
	p.Log("Fetching supported Purpur game versions...")

	// URL of the version manifest containing all Minecraft purpur versions
	const url = "https://api.purpurmc.org/v2/purpur"

	// Fetch and decode the purpur version manifest
	var versionData versionManifest
	if err := internal.FetchJSON(ctx, url, &versionData); err != nil {
		return nil, err
	}

	// Reverse the slice (higher versions first)
	versions := make([]string, 0, len(versionData.Versions))
	for i := len(versionData.Versions) - 1; i >= 0; i-- {
		versions = append(versions, versionData.Versions[i])
	}

	p.Log("Fetched %d Purpur game versions", len(versions))

	return versions, nil
}
