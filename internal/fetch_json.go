package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FetchJSON fetches JSON data from the given URL and decodes it into the provided variable.
//
// Parameters:
//   - url: the URL to fetch the JSON from.
//   - value: a pointer to the variable where the decoded JSON will be stored.
//
// Returns:
//   - error: an error if the HTTP request fails or the JSON cannot be decoded.
func FetchJSON(url string, value any) error {
	// Send HTTP GET request to the specified URL
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch JSON from %s: %w", url, err)
	}
	defer response.Body.Close()

	// Check for a successful HTTP response
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d when fetching JSON from %s", response.StatusCode, url)
	}

	// Decode the JSON response into the provided variable
	if err := json.NewDecoder(response.Body).Decode(value); err != nil {
		return fmt.Errorf("failed to decode JSON from %s: %w", url, err)
	}

	return nil
}
