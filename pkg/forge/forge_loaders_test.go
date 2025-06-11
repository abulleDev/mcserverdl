package forge_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/forge"
)

func TestForgeLoaders(t *testing.T) {
	const testVersion = "1.21.5"
	var latestFirst, oldestFirst []string

	t.Run("fetch latestFirst", func(t *testing.T) {
		var err error
		latestFirst, err = forge.Loaders(testVersion, true)
		if err != nil {
			t.Fatalf("Loaders(true) returned error: %v", err)
		}
		if len(latestFirst) == 0 {
			t.Fatal("Loaders(true) returned empty slice")
		}
	})

	t.Run("fetch oldestFirst", func(t *testing.T) {
		var err error
		oldestFirst, err = forge.Loaders(testVersion, false)
		if err != nil {
			t.Fatalf("Loaders(false) returned error: %v", err)
		}
		if len(oldestFirst) == 0 {
			t.Fatal("Loaders(false) returned empty slice")
		}
	})

	t.Run("check order", func(t *testing.T) {
		if len(latestFirst) != len(oldestFirst) {
			t.Fatalf("Length mismatch: latestFirst=%d, oldestFirst=%d", len(latestFirst), len(oldestFirst))
		}
		for i := range latestFirst {
			if latestFirst[i] != oldestFirst[len(oldestFirst)-1-i] {
				t.Fatalf("Order mismatch at index %d: latestFirst=%s, oldestFirst=%s", i, latestFirst[i], oldestFirst[len(oldestFirst)-1-i])
			}
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		_, err := forge.Loaders("invalid version", true)
		if err == nil {
			t.Error("expected error for invalid version, got nil")
		}
	})
}
