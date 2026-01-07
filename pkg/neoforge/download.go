package neoforge

import (
	"path/filepath"

	"github.com/abulleDev/mcserverdl/internal"
)

// Download downloads the NeoForge installer JAR to the specified installation directory.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.21.6", "25w14craftmine", "1.21").
//   - serverVersion: the NeoForge loader version.
//   - installDir: the directory where the installer JAR will be saved.
//   - onProgress: a callback function to report download progress.
//
// Returns:
//   - error: an error if the download fails.
func (p *Provider) Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	url, err := p.DownloadURL(gameVersion, serverVersion)
	if err != nil {
		return err
	}

	serverJarPath := filepath.Join(installDir, "installer.jar")
	return internal.Download(url, serverJarPath, onProgress)
}
