package vanilla

import (
	"context"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

type versionManifest struct {
	Versions []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"versions"`
}

// GameVersions fetches the list of all Minecraft vanilla versions from the official Mojang API version manifest.
// It uses a default background context.
func (p *Provider) GameVersions() ([]string, error) {
	return p.GameVersionsContext(context.Background())
}

// GameVersionsContext fetches the list of all Minecraft vanilla versions from the official Mojang API version manifest with context support.
//
// Parameters:
//   - ctx: the context to control the request lifetime.
//
// Returns:
//   - []string: a slice of Minecraft versions (e.g., "1.16.5", "15w14a", "1.18-pre2").
//   - error: an error if any HTTP or JSON decoding issues occur.
func (p *Provider) GameVersionsContext(ctx context.Context) ([]string, error) {
	p.Log("Fetching supported Vanilla game versions...")

	// URL of the version manifest containing all Minecraft vanilla versions
	const url = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

	// Fetch and decode the version manifest
	var versionData versionManifest
	if err := internal.FetchJSON(ctx, url, &versionData); err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(versionData.Versions))
	for _, version := range versionData.Versions {
		versions = append(versions, version.ID)
	}

	p.Log("Fetched %d vanilla game versions", len(versions))

	return versions, nil
}
