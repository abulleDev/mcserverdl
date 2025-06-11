package forge_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/forge"
)

func TestForgeVersions(t *testing.T) {
	var latestFirst, oldestFirst []string

	t.Run("fetch latestFirst", func(t *testing.T) {
		var err error
		latestFirst, err = forge.Versions(true)
		if err != nil {
			t.Fatalf("Versions(true) returned error: %v", err)
		}
		if len(latestFirst) == 0 {
			t.Fatal("Versions(true) returned empty slice")
		}
	})

	t.Run("fetch oldestFirst", func(t *testing.T) {
		var err error
		oldestFirst, err = forge.Versions(false)
		if err != nil {
			t.Fatalf("Versions(false) returned error: %v", err)
		}
		if len(oldestFirst) == 0 {
			t.Fatal("Versions(false) returned empty slice")
		}
	})

	t.Run("check order", func(t *testing.T) {
		if len(latestFirst) != len(oldestFirst) {
			t.Fatalf("Length mismatch: latestFirst=%d, oldestFirst=%d", len(latestFirst), len(oldestFirst))
		}
		for i := range latestFirst {
			if latestFirst[i] != oldestFirst[len(oldestFirst)-1-i] {
				t.Fatalf("Order mismatch at index %d: latestFirst=%q, oldestFirst=%q", i, latestFirst[i], oldestFirst[len(oldestFirst)-1-i])
			}
		}
	})
}
