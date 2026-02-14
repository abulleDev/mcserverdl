package provider_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/v2/pkg/provider"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/fabric"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/forge"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/neoforge"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/paper"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/purpur"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/vanilla"
)

func TestServerVersions(t *testing.T) {
	testCases := []struct {
		providerName       string
		gameVersion        string
		provider           provider.Provider
		expectFetchError   bool
		expectInvalidError bool
	}{
		{"Vanilla", "", vanilla.New(), true, true},
		{"Paper", "1.12.2", paper.New(), false, true},
		{"Fabric", "", fabric.New(), false, true},
		{"Forge", "1.21.5", forge.New(), false, true},
		{"NeoForge", "1.21.5", neoforge.New(), false, true},
		{"Purpur", "1.21.11", purpur.New(), false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.providerName, func(t *testing.T) {
			t.Parallel()
			t.Run("fetch support server versions", func(t *testing.T) {
				versions, err := tc.provider.ServerVersions(tc.gameVersion)
				if tc.expectFetchError {
					if err == nil {
						t.Fatal("expected error, got nil")
					}
				} else {
					if err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
					if len(versions) == 0 {
						t.Fatal("ServerVersions() returned empty slice")
					}
				}
			})

			t.Run("invalid version", func(t *testing.T) {
				_, err := tc.provider.ServerVersions("invalid version")
				if tc.expectInvalidError {
					if err == nil {
						t.Fatal("expected error, got nil")
					}
				} else {
					if err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
				}
			})
		})
	}
}
