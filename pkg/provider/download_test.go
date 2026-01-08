package provider_test

import (
	"os"
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/provider"
	"github.com/abulleDev/mcserverdl/pkg/provider/fabric"
	"github.com/abulleDev/mcserverdl/pkg/provider/forge"
	"github.com/abulleDev/mcserverdl/pkg/provider/neoforge"
	"github.com/abulleDev/mcserverdl/pkg/provider/paper"
	"github.com/abulleDev/mcserverdl/pkg/provider/vanilla"
)

func TestDownload(t *testing.T) {
	testCases := []struct {
		providerName  string
		gameVersion   string
		serverVersion string
		provider      provider.Provider
	}{
		{"Vanilla", "1.12.2", "", vanilla.New()},
		{"Paper", "1.12.2", "1620", paper.New()},
		{"Fabric", "1.21.5", "0.16.14", fabric.New()},
		{"Forge", "1.21.5", "55.0.23", forge.New()},
		{"NeoForge", "1.21.5", "21.5.75", neoforge.New()},
	}

	for _, tc := range testCases {
		t.Run(tc.providerName, func(t *testing.T) {
			installDir := t.TempDir()

			err := tc.provider.Download(tc.gameVersion, tc.serverVersion, installDir, nil)
			if err != nil {
				t.Fatalf("download failed: %v", err)
			}

			files, err := os.ReadDir(installDir)
			if err != nil {
				t.Fatalf("failed to read install dir: %v", err)
			}

			if len(files) == 0 {
				t.Error("expected downloaded files, but install dir is empty")
			}
		})
	}
}
