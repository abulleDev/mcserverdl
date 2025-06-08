package paper_test

import (
	"testing"

	"github.com/abulleDev/mcserverdl/pkg/paper"
)

func TestPaperDownloadURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, err := paper.DownloadURL("1.12.2", 1620)
		if err != nil {
			t.Fatalf("expected no error for valid version/build, got: %v", err)
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		_, err := paper.DownloadURL("invalid version", 0)
		if err == nil {
			t.Error("expected error for invalid version, got nil")
		}
	})

	t.Run("invalid build number", func(t *testing.T) {
		_, err := paper.DownloadURL("1.12.2", 0)
		if err == nil {
			t.Error("expected error for invalid build number, got nil")
		}
	})
}
