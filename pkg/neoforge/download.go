package neoforge

import (
	"path/filepath"

	"github.com/abulleDev/mcserverdl/internal"
)

func (p *Provider) Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	url, err := p.DownloadURL(gameVersion, serverVersion)
	if err != nil {
		return err
	}

	serverJarPath := filepath.Join(installDir, "installer.jar")
	return internal.Download(url, serverJarPath, onProgress)
}
