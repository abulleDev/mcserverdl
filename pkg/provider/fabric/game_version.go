package fabric

import "github.com/abulleDev/mcserverdl/internal"

type versionManifest []struct {
	Version string `json:"version"`
}

// GameVersions fetches the list of all Minecraft Fabric-supported game versions from the official FabricMC API.
//
// Returns:
//   - []string: a slice of Minecraft versions supported by Fabric (e.g., "1.20.5", "1.18-pre2", "20w51a").
//   - error: an error if any HTTP or JSON decoding issues occur.
func (p *Provider) GameVersions() ([]string, error) {
	p.Log("Fetching supported Fabric game versions...")

	// URL of the version manifest containing all Minecraft fabric versions
	const url = "https://meta2.fabricmc.net/v2/versions/game"

	// Fetch and decode the fabric version manifest
	var versionData versionManifest
	if err := internal.FetchJSON(url, &versionData); err != nil {
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
