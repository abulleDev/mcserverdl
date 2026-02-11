package paper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GameVersions fetches the list of all Minecraft paper server versions from the official PaperMC API version manifest.
// It uses a default background context.
func (p *Provider) GameVersions() ([]string, error) {
	return p.GameVersionsContext(context.Background())
}

// GameVersionsContext fetches the list of all Minecraft paper server versions from the official PaperMC API version manifest with context support.
//
// Parameters:
//   - ctx: the context to control the request lifetime.
//
// Returns:
//   - []string: a slice of Minecraft paper server versions (e.g., "1.16.5", "1.13-pre7").
//   - error: an error if any HTTP or JSON decoding issues occur.
func (p *Provider) GameVersionsContext(ctx context.Context) ([]string, error) {
	p.Log("Fetching supported Paper game versions...")

	// URL of the version manifest containing all Minecraft paper server versions
	const url = "https://fill.papermc.io/v3/projects/paper"

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	// Send HTTP GET request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON from %s: %w", url, err)
	}
	defer response.Body.Close()

	// Check for a successful HTTP response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d when fetching JSON from %s", response.StatusCode, url)
	}

	var versions []string

	decoder := json.NewDecoder(response.Body)

	// 1. Read the opening token of the top-level object ({).
	_, err = decoder.Token() // {
	if err != nil {
		return nil, fmt.Errorf("failed to read start of top-level JSON object: %w", err)
	}

	// 2. Iterate through the keys of the top-level object.
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("failed to read JSON key token: %w", err)
		}
		key := keyToken.(string)

		// 3. Find the "versions" key.
		if key == "versions" {
			// Read the opening token of the "versions" object ({).
			_, err := decoder.Token() // {
			if err != nil {
				return nil, fmt.Errorf("failed to read start of 'versions' object: %w", err)
			}

			// 4. Iterate inside the "versions" object to read the list of versions.
			for decoder.More() {
				// Read the version group key (e.g., "1.21"), but it's not used here.
				_, err := decoder.Token() // "1.21", "1.20", ...
				if err != nil {
					return nil, fmt.Errorf("failed to read version group key: %w", err)
				}

				// Decode the value for the key (the version array).
				var versionList []string
				if err := decoder.Decode(&versionList); err != nil {
					return nil, fmt.Errorf("failed to decode version array: %w", err)
				}

				// Append the read version list to the final slice.
				versions = append(versions, versionList...)
			}
			// Read the closing token of the "versions" object (}).
			_, err = decoder.Token() // }
			if err != nil {
				return nil, fmt.Errorf("failed to read end of 'versions' object: %w", err)
			}

		} else {
			// Skip the values of keys other than "versions".
			var ignoredValue any
			if err := decoder.Decode(&ignoredValue); err != nil {
				return nil, fmt.Errorf("failed to skip JSON value: %w", err)
			}
		}
	}

	p.Log("Fetched %d Paper game versions", len(versions))

	return versions, nil
}
