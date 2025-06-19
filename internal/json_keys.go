package internal

import (
	"encoding/json"
	"fmt"
	"io"
)

// ExtractJSONKeys extracts the top-level keys from a JSON object provided by an io.Reader.
//
// It expects the input to be a valid JSON object (i.e., starting with '{').
// The function iterates through all key-value pairs at the top level, collects the keys in the order
// they appear in the JSON string, and skips the values.
//
// Parameters:
//   - r: an io.Reader containing the JSON object.
//
// Returns:
//   - []string: a slice of strings containing the keys found at the top level of the JSON object, in the same order as in the input JSON.
//   - error: an error if the input is not a valid JSON object or if any decoding issues occur.
func ExtractJSONKeys(r io.Reader) ([]string, error) {
	// Create a new JSON decoder for the input reader
	decoder := json.NewDecoder(r)

	// Expect opening '{' token
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to read opening token: %w", err)
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '{' {
		return nil, fmt.Errorf("expected JSON object to start with '{', but got %v", token)
	}

	keys := make([]string, 0)

	// Iterate through all key-value pairs in the object
	for decoder.More() {
		// Read the key token
		token, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("failed to read key: %w", err)
		}
		key, ok := token.(string)
		if !ok {
			return nil, fmt.Errorf("expected a string key, but got a token of type %T", token)
		}
		keys = append(keys, key)

		// Skip the value by decoding into RawMessage
		var dummy json.RawMessage
		if err := decoder.Decode(&dummy); err != nil {
			return nil, fmt.Errorf("failed to skip value for key %q: %w", key, err)
		}
	}

	return keys, nil
}
