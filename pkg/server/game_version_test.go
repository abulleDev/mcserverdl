package server_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/fabric"
	"github.com/abulleDev/mcserverdl/pkg/forge"
	"github.com/abulleDev/mcserverdl/pkg/neoforge"
	"github.com/abulleDev/mcserverdl/pkg/paper"
	"github.com/abulleDev/mcserverdl/pkg/server"
	"github.com/abulleDev/mcserverdl/pkg/vanilla"
)

func TestGameVersions(t *testing.T) {
	testCases := []struct {
		providerName     string
		provider         server.Provider
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
