package fabric

import (
	"context"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

type versionManifest []struct {
	Version string `json:"version"`
}

// GameVersions fetches the list of all Minecraft Fabric-supported game versions from the official FabricMC API.
// It uses a default background context.
func (p *Provider) GameVersions() ([]string, error) {
	return p.GameVersionsContext(context.Background())
}

// GameVersionsContext fetches the list of all Minecraft Fabric-supported game versions from the official FabricMC API with context support.
//
// Parameters:
//   - ctx: the context to control the request lifetime.
//
// Returns:
//   - []string: a slice of Minecraft versions supported by Fabric (e.g., "1.20.5", "1.18-pre2", "20w51a").
//   - error: an error if any HTTP or JSON decoding issues occur.
func (p *Provider) GameVersionsContext(ctx context.Context) ([]string, error) {
	p.Log("Fetching supported Fabric game versions...")

	// URL of the version manifest containing all Minecraft fabric versions
	const url = "https://meta2.fabricmc.net/v2/versions/game"

	// Fetch and decode the fabric version manifest
	var versionData versionManifest
	if err := internal.FetchJSON(ctx, url, &versionData); err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(versionData))
	// Build the slice from first to last (higher versions first)
	for _, version := range versionData {
		versions = append(versions, version.Version)
	}

	p.Log("Fetched %d Fabric game versions", len(versions))

	return versions, nil
}
