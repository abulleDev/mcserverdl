package fabric

import (
	"path/filepath"

	"github.com/abulleDev/mcserverdl/internal"
)

// Download downloads the Fabric server JAR to the specified installation directory.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.20.5", "1.18-pre2", "20w51a").
//   - serverVersion: the Fabric loader version.
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
