package neoforge

import (
	"fmt"
	"strings"
)

func parseGameVersion(loaderVersion string) (string, error) {
	versionParts := strings.SplitN(loaderVersion, ".", 4)

	switch versionParts[0] {
	case "0":
		// Extract game version from legacy snapshot format
		// Format:
		// 	- 0
		// 	- <game snapshot version>
		// 	- <neoforge build>
		// Example: "0.25w14craftmine.5-beta" -> "25w14craftmine"
		if len(versionParts) != 3 {
			return "", fmt.Errorf("invalid legacy snapshot version format: %s", loaderVersion)
		}

		return versionParts[1], nil
	case "20", "21":
		// Extract game version from legacy release format
		// Format:
		// 	- <game major version>
		// 	- <game minor version>
		// 	- <neoforge build>
		// Example: "21.0.142-beta" -> "1.21"
		if len(versionParts) != 3 {
			return "", fmt.Errorf("invalid legacy release version format: %s", loaderVersion)
		}

		return strings.TrimSuffix(fmt.Sprintf("1.%s.%s", versionParts[0], versionParts[1]), ".0"), nil
	default:
		// Extract game version from new format (Release & Snapshot)
		// Format:
		// 	- <published year>
		// 	- <game major version>
		// 	- <game patch version>
		// 	- <neoforge build> [+<game snapshot version>]
		// Examples:
		// 	- Release:  "26.1.1.1-beta" -> "21.1.1"
		// 	- Snapshot: "26.1.0.0-alpha.14+snapshot-11" -> "21.1-snapshot-11"
		if len(versionParts) != 4 {
			return "", fmt.Errorf("invalid version format: %s", loaderVersion)
		}

		gameVersion := strings.TrimSuffix(strings.Join(versionParts[0:3], "."), ".0")

		plusIndex := strings.Index(versionParts[3], "+")
		if plusIndex == -1 {
			return gameVersion, nil
		} else {
			return gameVersion + "-" + versionParts[3][plusIndex+1:], nil
		}
	}
}
