package vanilla

import (
	"path/filepath"

	"github.com/abulleDev/mcserverdl/v2/internal"
)

// Download downloads the vanilla server JAR to the specified installation directory.
//
// Parameters:
//   - gameVersion: the Minecraft version string (e.g., "1.16.5", "15w14a", "1.18-pre2").
//   - serverVersion: ignored for vanilla.
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

	p.Log("Downloading server...")

	serverJarPath := filepath.Join(installDir, "server.jar")
	err = internal.Download(url, serverJarPath, onProgress)

	p.Log("Successfully downloaded server to %s", installDir)

	return err
}
