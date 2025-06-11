package forge_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/forge"
)

func TestForgeDownloadURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, err := forge.DownloadURL("1.21.5", "55.0.23")
		if err != nil {
			t.Fatalf("expected no error for valid version/loader, got: %v", err)
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		_, err := forge.DownloadURL("invalid version", "0")
		if err == nil {
			t.Error("expected error for invalid version, got nil")
		}
	})

	t.Run("invalid loader version", func(t *testing.T) {
		_, err := forge.DownloadURL("1.21.5", "0")
		if err == nil {
			t.Error("expected error for invalid loader version, got nil")
		}
	})
}
