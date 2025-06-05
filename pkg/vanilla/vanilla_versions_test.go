package vanilla_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/vanilla"
)

func TestVanillaVersions(t *testing.T) {
	latestFirst, err := vanilla.Versions(true)
	if err != nil {
		t.Fatalf("Versions(true) returned error: %v", err)
	}
	if len(latestFirst) == 0 {
		t.Fatal("Versions(true) returned empty slice")
	}

	oldestFirst, err := vanilla.Versions(false)
	if err != nil {
		t.Fatalf("Versions(false) returned error: %v", err)
	}
	if len(oldestFirst) == 0 {
		t.Fatal("Versions(false) returned empty slice")
	}

	if len(latestFirst) != len(oldestFirst) {
		t.Fatalf("Length mismatch: latestFirst=%d, oldestFirst=%d", len(latestFirst), len(oldestFirst))
	}
	for i := range latestFirst {
		if latestFirst[i] != oldestFirst[len(oldestFirst)-1-i] {
			t.Errorf("Order mismatch at index %d: latestFirst=%q, oldestFirst=%q", i, latestFirst[i], oldestFirst[len(oldestFirst)-1-i])
		}
	}

	latest := latestFirst[0]
	lastOld := oldestFirst[len(oldestFirst)-1]
	if latest != lastOld {
		t.Errorf(
			"Mismatch between sorting orders:\n latestFirst[0] = %q\n oldestFirst[last] = %q",
			latest, lastOld,
		)
	}
}
