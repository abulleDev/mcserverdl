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

func TestDownloadURL(t *testing.T) {
	testCases := []struct {
		providerName      string
		gameVersion       string
		serverVersion     string
		provider          provider.Provider
		expectFetchError  bool
		expectGameError   bool
		expectServerError bool
	}{
		{"Vanilla", "1.12.2", "", vanilla.New(), false, true, false},
		{"Paper", "1.12.2", "1620", paper.New(), false, true, true},
		{"Fabric", "1.21.5", "0.16.14", fabric.New(), false, true, true},
		{"Forge", "1.21.5", "55.0.23", forge.New(), false, true, true},
		{"NeoForge", "1.21.5", "21.5.75", neoforge.New(), false, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.providerName, func(t *testing.T) {
			t.Run("fetch server download url of valid version", func(t *testing.T) {
				_, err := tc.provider.DownloadURL(tc.gameVersion, tc.serverVersion)
				if tc.expectFetchError {
					if err == nil {
						t.Fatal("expected error, got nil")
					}
				} else {
					if err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
				}
			})

			t.Run("invalid game version", func(t *testing.T) {
				_, err := tc.provider.DownloadURL("invalid game version", tc.serverVersion)
				if tc.expectGameError {
					if err == nil {
						t.Fatal("expected error, got nil")
					}
				} else {
					if err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
				}
			})

			t.Run("invalid server version", func(t *testing.T) {
				_, err := tc.provider.DownloadURL(tc.gameVersion, "invalid server version")
				if tc.expectServerError {
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
