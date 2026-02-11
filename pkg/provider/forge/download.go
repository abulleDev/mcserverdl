package forge

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abulleDev/mcserverdl/v2/internal"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/vanilla"
)

// Download downloads the Forge server files to the specified installation directory.
// It uses a default background context.
func (p *Provider) Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	return p.DownloadContext(context.Background(), gameVersion, serverVersion, installDir, onProgress)
}

// DownloadContext downloads the Forge server files to the specified installation directory with context support.
// It handles both standard installer JARs and legacy zip patches (which require merging with a vanilla server JAR).
//
// Parameters:
//   - ctx: the context to control the download cancellation.
//   - gameVersion: the Minecraft version string (e.g., "1.21.6", "1.7.10-pre4", "1.4").
//   - serverVersion: the Forge loader version.
//   - installDir: the directory where the server files will be saved.
//   - onProgress: a callback function to report download progress.
//
// Returns:
//   - error: an error if the download or patching (for legacy versions) fails.
func (p *Provider) DownloadContext(ctx context.Context, gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error {
	url, err := p.DownloadURLContext(ctx, gameVersion, serverVersion)
	if err != nil {
		return err
	}

	if strings.HasSuffix(url, ".jar") {
		// Case 1: The URL points to a standard installer JAR
		p.Log("Downloading Forge installer...")
		installerPath := filepath.Join(installDir, "installer.jar")
		if err := internal.Download(ctx, url, installerPath, onProgress); err != nil {
			return err
		}
		p.Log("Installer downloaded. Please run the following command in the installation directory to complete the server setup:")
		p.Log("java -jar installer.jar --installServer")
	} else if strings.HasSuffix(url, ".zip") {
		// Case 2: The URL points to a patch file that needs to be applied to a vanilla server
		patchPath := filepath.Join(installDir, "patch.zip")
		vanillaPath := filepath.Join(installDir, "vanilla.jar")
		finalJarPath := filepath.Join(installDir, "server.jar")

		// Clean up temporary files used during the patching process
		defer func() {
			os.Remove(patchPath)
			os.Remove(vanillaPath)
		}()

		// Download the patch file
		p.Log("Downloading Forge patch file...")
		if err := internal.Download(ctx, url, patchPath, onProgress); err != nil {
			return err
		}
		p.Log("Download complete!")

		// Download the corresponding vanilla server JAR
		p.Log("Downloading vanilla server for %s...", gameVersion)
		vanillaURL, err := vanilla.New().DownloadURLContext(ctx, gameVersion, "")
		if err != nil {
			return err
		}
		if err := internal.Download(ctx, vanillaURL, vanillaPath, onProgress); err != nil {
			return err
		}
		p.Log("Download complete!")

		// Patch the server
		p.Log("Patching vanilla server...")
		if err := internal.MergeZips(ctx, vanillaPath, patchPath, finalJarPath); err != nil {
			return err
		}
		p.Log("Successfully created Forge server to %s", installDir)
	} else {
		return fmt.Errorf("unexpected URL format: %s", url)
	}

	return nil
}
