package paper

import (
	"path/filepath"

	"github.com/abulleDev/mcserverdl/internal"
)

// Download downloads the PaperMC server JAR to the specified installation directory.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.16.5", "1.13-pre7").
//   - serverVersion: the Paper build number.
//   - installDir: the directory where the server JAR will be saved.
//   - onProgress: a callback function to report download progress.
//
// Returns:
//   - error: an error if the download fails.
func (p *Provider) Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	url, err := p.DownloadURL(gameVersion, serverVersion)
	if err != nil {
		return err
	}

	serverJarPath := filepath.Join(installDir, "server.jar")
	return internal.Download(url, serverJarPath, onProgress)
}
