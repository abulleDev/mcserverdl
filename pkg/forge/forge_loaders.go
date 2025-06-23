package forge

import (
	"fmt"
	"strings"

	"github.com/abulleDev/mcserverdl/internal"
)

// Loaders fetches a list of available Forge loader versions for a given Minecraft version.
// It retrieves the data from the official Forge maven metadata.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.21.6", "1.7.10-pre4", "1.4").
//   - latestFirst: if true, returns the loader versions with higher versions first. If false, returns the loader versions with lower versions first.
//
// Returns:
//   - []string: a slice of Forge loader versions (e.g., "56.0.3", "14.23.4.2720").
//   - error: an error if the game version is not supported or if any HTTP or JSON decoding issues occur.
func Loaders(gameVersion string, latestFirst bool) ([]string, error) {
	// URL of the version manifest containing all Minecraft forge versions
	const url = "https://files.minecraftforge.net/net/minecraftforge/forge/maven-metadata.json"

	// Fetch and decode the forge loader manifest
	var loaderData map[string][]string
	if err := internal.FetchJSON(url, &loaderData); err != nil {
		return nil, err
	}

	// Some game versions have different naming conventions in the manifest
	var forgeStyleVersion string
	switch gameVersion {
	case "1.7.10-pre4":
		forgeStyleVersion = "1.7.10_pre4"
	case "1.4":
		forgeStyleVersion = "1.4.0"
	default:
		forgeStyleVersion = gameVersion
	}

	// Raw loader versions from the manifest (e.g., "1.7.10-10.13.3.1401-1710ls")
	rawLoaderVersions, ok := loaderData[forgeStyleVersion]
	if !ok {
		return nil, fmt.Errorf("unsupported game version: %s", gameVersion)
	}

	// Refined loader versions (e.g., "10.13.3.1401")
	refinedLoadersVersions := make([]string, 0, len(rawLoaderVersions))
	if latestFirst {
		// Build the slice from last to first (higher versions first)
		for i := len(rawLoaderVersions) - 1; i >= 0; i-- {
			refinedLoadersVersions = append(refinedLoadersVersions, strings.Split(rawLoaderVersions[i], "-")[1])
		}
	} else {
		// Build the slice from first to last (lower versions first)
		for _, rawLoaderVersion := range rawLoaderVersions {
			refinedLoadersVersions = append(refinedLoadersVersions, strings.Split(rawLoaderVersion, "-")[1])
		}
	}

	return refinedLoadersVersions, nil
}
