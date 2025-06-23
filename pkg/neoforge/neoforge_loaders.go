package neoforge

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

// Loaders fetches a list of available NeoForge loader versions for a given Minecraft version.
// It retrieves the data from the official NeoForged maven metadata.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.21.6", "25w14craftmine", "1.21").
//   - latestFirst: if true, returns the loader versions with higher versions first. If false, returns the loader versions with lower versions first.
//
// Returns:
//   - []string: a slice of NeoForge loader versions (e.g., "21.0.142-beta", "0.25w14craftmine.5-beta").
//   - error: an error if the game version is not supported or if any HTTP or XML decoding issues occur.
func Loaders(gameVersion string, latestFirst bool) ([]string, error) {
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

	// Filter loader versions that match the requested game version
	matchingLoaderVersions := make([]string, 0, len(versionData.Versioning.Versions.Version))
	for _, loaderVersion := range versionData.Versioning.Versions.Version {
		isSnapshot := strings.HasPrefix(loaderVersion, "0.")

		var currentGameVersion string
		if isSnapshot {
			// Extract game version from snapshot format (e.g., "0.25w14craftmine.5-beta" -> "25w14craftmine")
			firstDotIndex := strings.Index(loaderVersion, ".")
			lastDotIndex := strings.LastIndex(loaderVersion, ".")
			if firstDotIndex == -1 || lastDotIndex == -1 || firstDotIndex == lastDotIndex {
				return nil, fmt.Errorf("invalid snapshot version format: %s", loaderVersion)
			}
			currentGameVersion = loaderVersion[firstDotIndex+1 : lastDotIndex]
		} else {
			// Extract game version from release format (e.g., "21.0.142-beta" -> "1.21")
			upToLastDot := loaderVersion[:strings.LastIndex(loaderVersion, ".")]
			currentGameVersion = "1." + strings.TrimSuffix(upToLastDot, ".0")
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

	// Return the versions as-is (lower versions first)
	if !latestFirst {
		return matchingLoaderVersions, nil
	}

	// Reverse the slice (higher versions first)
	for i, j := 0, len(matchingLoaderVersions)-1; i < j; i, j = i+1, j-1 {
		matchingLoaderVersions[i], matchingLoaderVersions[j] = matchingLoaderVersions[j], matchingLoaderVersions[i]
	}
	return matchingLoaderVersions, nil
}
