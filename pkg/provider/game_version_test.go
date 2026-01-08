package provider_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/provider"
	"github.com/abulleDev/mcserverdl/pkg/provider/fabric"
	"github.com/abulleDev/mcserverdl/pkg/provider/forge"
	"github.com/abulleDev/mcserverdl/pkg/provider/neoforge"
	"github.com/abulleDev/mcserverdl/pkg/provider/paper"
	"github.com/abulleDev/mcserverdl/pkg/provider/vanilla"
)

func TestGameVersions(t *testing.T) {
	testCases := []struct {
		providerName     string
		provider         provider.Provider
		expectFetchError bool
	}{
		{"Vanilla", vanilla.New(), false},
		{"Paper", paper.New(), false},
		{"Fabric", fabric.New(), false},
		{"Forge", forge.New(), false},
		{"NeoForge", neoforge.New(), false},
	}

	for _, tc := range testCases {
		t.Run(tc.providerName, func(t *testing.T) {
			t.Parallel()
			t.Run("fetch support game versions", func(t *testing.T) {
				versions, err := tc.provider.GameVersions()
				if tc.expectFetchError {
					if err == nil {
						t.Fatal("expected error, got nil")
					}
				} else {
					if err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
					if len(versions) == 0 {
						t.Fatal("GameVersions() returned empty slice")
					}
				}
			})
		})
	}
}
