package purpur

import (
	"context"
	"path/filepath"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

// Download downloads the PurpurMC server JAR to the specified installation directory.
// It uses a default background context.
func (p *Provider) Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	return p.DownloadContext(context.Background(), gameVersion, serverVersion, installDir, onProgress)
}

// DownloadContext downloads the PurpurMC server JAR to the specified installation directory with context support.
//
// Parameters:
//   - ctx: the context to control the download cancellation.
//   - gameVersion: the Minecraft version string (e.g., "1.21.11", "1.14.1").
//   - serverVersion: the Purpur build number.
//   - installDir: the directory where the server JAR will be saved.
//   - onProgress: a callback function to report download progress.
//
// Returns:
//   - error: an error if the download fails.
func (p *Provider) DownloadContext(ctx context.Context, gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	url, err := p.DownloadURLContext(ctx, gameVersion, serverVersion)
	if err != nil {
		return err
	}

	p.Log("Downloading server...")

	serverJarPath := filepath.Join(installDir, "server.jar")
	err = internal.Download(ctx, url, serverJarPath, onProgress)

	p.Log("Successfully downloaded server to %s", installDir)

	return err
}
