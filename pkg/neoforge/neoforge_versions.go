package neoforge

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

type neoforgeVersionManifest struct {
	Versioning struct {
		Versions struct {
			Version []string `xml:"version"`
		} `xml:"versions"`
	} `xml:"versioning"`
}

// Versions fetches the list of all Minecraft NeoForge-supported game versions from the official NeoForged maven metadata.
//
// Parameters:
//   - latestFirst: if true, returns the versions with higher versions first. If false, returns the versions with lower versions first.
//
// Returns:
//   - []string: a slice of Minecraft versions supported by NeoForge (e.g., "1.21.6", "25w14craftmine", "1.21").
//   - error: an error if any HTTP or XML decoding issues occur.
func Versions(latestFirst bool) ([]string, error) {
	// URL of the version manifest containing all Minecraft neoforge versions
	const url = "https://maven.neoforged.net/releases/net/neoforged/neoforge/maven-metadata.xml"

	// Send HTTP GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch XML from %s: %w", url, err)
	}
	defer response.Body.Close()

	// Check for a successful HTTP response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d when fetching XML from %s", response.StatusCode, url)
	}

	// Decode the XML response into the provided variable
	var versionData neoforgeVersionManifest
	if err := xml.NewDecoder(response.Body).Decode(&versionData); err != nil {
		return nil, fmt.Errorf("failed to decode XML from %s: %w", url, err)
	}

	// Use a map to store unique game versions to avoid duplicates
	gameVersions := make([]string, 0, len(versionData.Versioning.Versions.Version))
	versionSet := map[string]struct{}{}

	// Iterate over all loader versions to extract the corresponding game version
	for _, loaderVersion := range versionData.Versioning.Versions.Version {
		isSnapshot := strings.HasPrefix(loaderVersion, "0.")

		var gameVersion string
		if isSnapshot {
			// Extract game version from snapshot format (e.g., "0.25w14craftmine.5-beta" -> "25w14craftmine")
			firstDotIndex := strings.Index(loaderVersion, ".")
			lastDotIndex := strings.LastIndex(loaderVersion, ".")
			if firstDotIndex == -1 || lastDotIndex == -1 || firstDotIndex == lastDotIndex {
				return nil, fmt.Errorf("invalid snapshot version format: %s", loaderVersion)
			}
			gameVersion = loaderVersion[firstDotIndex+1 : lastDotIndex]
		} else {
			// Extract game version from release format (e.g., "21.0.142-beta" -> "1.21")
			upToLastDot := loaderVersion[:strings.LastIndex(loaderVersion, ".")]
			gameVersion = "1." + strings.TrimSuffix(upToLastDot, ".0")
		}

		// Add the game version if it's new
		if _, exist := versionSet[gameVersion]; !exist {
			versionSet[gameVersion] = struct{}{}
			gameVersions = append(gameVersions, gameVersion)
		}
	}

	// Return the versions as-is (lower versions first)
	if !latestFirst {
		return gameVersions, nil
	}

	// Reverse the slice (higher versions first)
	for i, j := 0, len(gameVersions)-1; i < j; i, j = i+1, j-1 {
		gameVersions[i], gameVersions[j] = gameVersions[j], gameVersions[i]
	}
	return gameVersions, nil
}
