package neoforge

import (
	"context"
	"path/filepath"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

// Download downloads the NeoForge installer JAR to the specified installation directory.
// It uses a default background context.
func (p *Provider) Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	return p.DownloadContext(context.Background(), gameVersion, serverVersion, installDir, onProgress)
}

// DownloadContext downloads the NeoForge installer JAR to the specified installation directory with context support.
//
// Parameters:
//   - ctx: the context to control the download cancellation.
//   - gameVersion: the Minecraft version string (e.g., "1.21.6", "25w14craftmine", "1.21").
//   - serverVersion: the NeoForge loader version.
//   - installDir: the directory where the installer JAR will be saved.
//   - onProgress: a callback function to report download progress.
//
// Returns:
//   - error: an error if the download fails.
func (p *Provider) DownloadContext(ctx context.Context, gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	url, err := p.DownloadURLContext(ctx, gameVersion, serverVersion)
	if err != nil {
		return err
	}

	p.Log("Downloading NeoForge installer...")

	serverJarPath := filepath.Join(installDir, "installer.jar")
	err = internal.Download(ctx, url, serverJarPath, onProgress)

	p.Log("Installer downloaded. Please run the following command in the installation directory to complete the server setup:")
	p.Log("java -jar installer.jar --installServer")

	return err
}
