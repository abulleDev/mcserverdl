package neoforge

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

type versionManifest struct {
	Versioning struct {
		Versions struct {
			Version []string `xml:"version"`
		} `xml:"versions"`
	} `xml:"versioning"`
}

// GameVersions fetches the list of all Minecraft NeoForge-supported game versions from the official NeoForged maven metadata.
// It uses a default background context.
func (p *Provider) GameVersions() ([]string, error) {
	return p.GameVersionsContext(context.Background())
}

// GameVersionsContext fetches the list of all Minecraft NeoForge-supported game versions from the official NeoForged maven metadata with context support.
//
// Parameters:
//   - ctx: the context to control the request lifetime.
//
// Returns:
//   - []string: a slice of Minecraft versions supported by NeoForge (e.g., "1.21.6", "25w14craftmine", "1.21").
//   - error: an error if any HTTP or XML decoding issues occur.
func (p *Provider) GameVersionsContext(ctx context.Context) ([]string, error) {
	p.Log("Fetching supported NeoForge game versions...")

	// URL of the version manifest containing all Minecraft neoforge versions
	const url = "https://maven.neoforged.net/releases/net/neoforged/neoforge/maven-metadata.xml"

	p.Log("Fetching NeoForge-supported game versions...")

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	// Send HTTP GET request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch XML from %s: %w", url, err)
	}
	defer response.Body.Close()

	// Check for a successful HTTP response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d when fetching XML from %s", response.StatusCode, url)
	}

	// Decode the XML response into the provided variable
	var versionData versionManifest
	if err := xml.NewDecoder(response.Body).Decode(&versionData); err != nil {
		return nil, fmt.Errorf("failed to decode XML from %s: %w", url, err)
	}

	// Use a map to store unique game versions to avoid duplicates
	gameVersions := make([]string, 0, len(versionData.Versioning.Versions.Version))
	versionSet := map[string]struct{}{}

	// Iterate over all loader versions to extract the corresponding game version
	for _, loaderVersion := range versionData.Versioning.Versions.Version {
		gameVersion, err := parseGameVersion(loaderVersion)
		if err != nil {
			return nil, err
		}

		// Add the game version if it's new
		if _, exist := versionSet[gameVersion]; !exist {
			versionSet[gameVersion] = struct{}{}
			gameVersions = append(gameVersions, gameVersion)
		}
	}

	// Reverse the slice (higher versions first)
	for i, j := 0, len(gameVersions)-1; i < j; i, j = i+1, j-1 {
		gameVersions[i], gameVersions[j] = gameVersions[j], gameVersions[i]
	}

	p.Log("Fetched %d NeoForge game versions", len(gameVersions))
	return gameVersions, nil
}
