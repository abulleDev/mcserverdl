package forge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abulleDev/mcserverdl/internal"
	"github.com/abulleDev/mcserverdl/pkg/vanilla"
)

func (p *Provider) Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	url, err := p.DownloadURL(gameVersion, serverVersion)
	if err != nil {
		return err
	}

	if strings.HasSuffix(url, ".jar") {
		// Case 1: The URL points to a standard installer JAR
		installerPath := filepath.Join(installDir, "installer.jar")
		if err := internal.Download(url, installerPath, onProgress); err != nil {
			return err
		}
	} else if strings.HasSuffix(url, ".zip") {
		// Case 2: The URL points to a patch file that needs to be applied to a vanilla server
		patchPath := filepath.Join(installDir, "patch.zip")
		vanillaPath := filepath.Join(installDir, "vanilla.jar")
		finalJarPath := filepath.Join(installDir, "server.jar")

		// Download the patch file
		if err := internal.Download(url, patchPath, onProgress); err != nil {
			return err
		}

		// Download the corresponding vanilla server JAR
		vanillaURL, err := vanilla.New().DownloadURL(gameVersion, "")
		if err != nil {
			return err
		}
		if err := internal.Download(vanillaURL, vanillaPath, onProgress); err != nil {
			return err
		}

		// Patch the server
		if err := internal.MergeZips(vanillaPath, patchPath, finalJarPath); err != nil {
			return err
		}

		// Clean path files
		defer os.Remove(patchPath)
		defer os.Remove(vanillaPath)
	} else {
		return fmt.Errorf("unexpected URL format: %s", url)
	}

	return nil
}
