package paper_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/paper"
)

func TestPaperBuilds(t *testing.T) {
	const testVersion = "1.12.2"
	var latestFirst, oldestFirst []int

	t.Run("fetch latestFirst", func(t *testing.T) {
		var err error
		latestFirst, err = paper.Builds(testVersion, true)
		if err != nil {
			t.Fatalf("Builds(true) returned error: %v", err)
		}
		if len(latestFirst) == 0 {
			t.Fatal("Builds(true) returned empty slice")
		}
	})

	t.Run("fetch oldestFirst", func(t *testing.T) {
		var err error
		oldestFirst, err = paper.Builds(testVersion, false)
		if err != nil {
			t.Fatalf("Builds(false) returned error: %v", err)
		}
		if len(oldestFirst) == 0 {
			t.Fatal("Builds(false) returned empty slice")
		}
	})

	t.Run("check order", func(t *testing.T) {
		if len(latestFirst) != len(oldestFirst) {
			t.Fatalf("Length mismatch: latestFirst=%d, oldestFirst=%d", len(latestFirst), len(oldestFirst))
		}
		for i := range latestFirst {
			if latestFirst[i] != oldestFirst[len(oldestFirst)-1-i] {
				t.Fatalf("Order mismatch at index %d: latestFirst=%d, oldestFirst=%d", i, latestFirst[i], oldestFirst[len(oldestFirst)-1-i])
			}
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		_, err := paper.Builds("invalid version", true)
		if err == nil {
			t.Error("expected error for invalid version, got nil")
		}
	})
}
