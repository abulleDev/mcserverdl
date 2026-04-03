package neoforge

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// ServerVersions fetches a list of available NeoForge loader versions for a given Minecraft version.
// It uses a default background context.
func (p *Provider) ServerVersions(gameVersion string) ([]string, error) {
	return p.ServerVersionsContext(context.Background(), gameVersion)
}

// ServerVersionsContext fetches a list of available NeoForge loader versions for a given Minecraft version with context support.
// It retrieves the data from the official NeoForged maven metadata.
//
// Parameters:
//   - ctx: the context to control the request lifetime.
//   - gameVersion: the Minecraft version string (e.g., "1.21.6", "25w14craftmine", "1.21").
//
// Returns:
//   - []string: a slice of NeoForge loader versions (e.g., "21.0.142-beta", "0.25w14craftmine.5-beta").
//   - error: an error if the game version is not supported or if any HTTP or XML decoding issues occur.
func (p *Provider) ServerVersionsContext(ctx context.Context, gameVersion string) ([]string, error) {
	p.Log("Fetching NeoForge server versions (loaders) for %s...", gameVersion)

	// URL of the version manifest containing all Minecraft neoforge versions
	const url = "https://maven.neoforged.net/releases/net/neoforged/neoforge/maven-metadata.xml"

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

	// Filter loader versions that match the requested game version
	matchingLoaderVersions := make([]string, 0, len(versionData.Versioning.Versions.Version))
	for _, loaderVersion := range versionData.Versioning.Versions.Version {
		currentGameVersion, err := parseGameVersion(loaderVersion)
		if err != nil {
			return nil, err
		}

		// If the extracted game version matches the requested one, add the raw loader version to our list
		if currentGameVersion == gameVersion {
			matchingLoaderVersions = append(matchingLoaderVersions, loaderVersion)
		}
	}

	// If no matching loaders were found, the game version is unsupported
	if len(matchingLoaderVersions) == 0 {
		return nil, fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Reverse the slice (higher versions first)
	for i, j := 0, len(matchingLoaderVersions)-1; i < j; i, j = i+1, j-1 {
		matchingLoaderVersions[i], matchingLoaderVersions[j] = matchingLoaderVersions[j], matchingLoaderVersions[i]
	}

	p.Log("Fetched %d NeoForge loader versions for %s", len(matchingLoaderVersions), gameVersion)
	return matchingLoaderVersions, nil
}
