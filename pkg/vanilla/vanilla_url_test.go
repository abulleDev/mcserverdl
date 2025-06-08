package vanilla_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/vanilla"
)

func TestDownloadURL(t *testing.T) {
	t.Run("valid version", func(t *testing.T) {
		_, err := vanilla.DownloadURL("1.12.2")
		if err != nil {
			t.Fatalf("expected no error for valid version, got: %v", err)
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		_, err := vanilla.DownloadURL("invalid version")
		if err == nil {
			t.Error("expected error for invalid version, got nil")
		}
	})
}
