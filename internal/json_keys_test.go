package internal

import (
	"strings"
	"testing"
)

func TestExtractJSONKeys(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const validJSONString = `{"one":null,"two":null,"three":null,"four":null,"five":null}`
		result, err := ExtractJSONKeys(strings.NewReader(validJSONString))
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		want := []string{"one", "two", "three", "four", "five"}
		if len(result) != len(want) {
			t.Fatalf("expected %d keys, got %d", len(want), len(result))
		}
		for i, k := range want {
			if result[i] != k {
				t.Errorf("expected key %q at index %d, got %q", k, i, result[i])
			}
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		const invalidJSONString = `{"error":invalid value}`
		_, err := ExtractJSONKeys(strings.NewReader(invalidJSONString))
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})

	t.Run("not an object", func(t *testing.T) {
		const arrayJSONString = `[1, 2, 3, 4, 5]`
		_, err := ExtractJSONKeys(strings.NewReader(arrayJSONString))
		if err == nil {
			t.Error("expected error for non-object JSON, got nil")
		}
	})

	t.Run("empty object", func(t *testing.T) {
		const emptyObject = `{}`
		result, err := ExtractJSONKeys(strings.NewReader(emptyObject))
		if err != nil {
			t.Fatalf("expected no error for empty object, got: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected 0 keys for empty object, got %d", len(result))
		}
	})
}
