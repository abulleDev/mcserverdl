package provider

import "context"

// Logger defines the interface for logging messages.
// It abstraction allows the application to inject standard log.Logger,
// or any custom logger that supports Printf.
type Logger interface {
	Printf(format string, v ...any)
}

// BaseProvider implements the logging functionality shared by all providers.
// It is intended to be embedded in specific provider implementations.
// Since the logger can be nil, methods check for its existence before logging.
type BaseProvider struct {
	logger Logger
}

// SetLogger sets the logger instance for the provider.
// This allows external callers (like main) to inject a logger.
func (b *BaseProvider) SetLogger(l Logger) {
	b.logger = l
}

// Log prints a formatted message if a logger is configured.
// It is a helper to avoid repeated nil checks in the provider code.
func (b *BaseProvider) Log(format string, v ...any) {
	if b.logger == nil {
		return
	}

	b.logger.Printf(format, v...)
}

// Provider defines the standard interface that all Minecraft server providers must implement.
type Provider interface {
	// GameVersions returns a list of available game versions (e.g., "1.16.5", "15w14a", "1.18-pre2").
	// It is equivalent to calling GameVersionsContext with context.Background().
	GameVersions() ([]string, error)

	// GameVersionsContext returns a list of available game versions with context support.
	GameVersionsContext(ctx context.Context) ([]string, error)

	// ServerVersions returns a list of available server builds/loader versions for a specific game version.
	// It is equivalent to calling ServerVersionsContext with context.Background().
	ServerVersions(gameVersion string) ([]string, error)

	// ServerVersionsContext returns a list of available server builds/loader versions with context support.
	ServerVersionsContext(ctx context.Context, gameVersion string) ([]string, error)

	// DownloadURL returns the direct download URL for the server jar.
	// It is equivalent to calling DownloadURLContext with context.Background().
	DownloadURL(gameVersion, serverVersion string) (string, error)

	// DownloadURLContext returns the direct download URL for the server jar with context support.
	DownloadURLContext(ctx context.Context, gameVersion, serverVersion string) (string, error)

	// Download downloads the server jar to the specified directory.
	// onProgress is called periodically with bytes downloaded and total file size.
	// It is equivalent to calling DownloadContext with context.Background().
	Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error

	// DownloadContext downloads the server jar to the specified directory with context support.
	// onProgress is called periodically with bytes downloaded and total file size.
	DownloadContext(ctx context.Context, gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error

	// SetLogger injects a logger into the provider.
	SetLogger(l Logger)

	// Log logs a message using the injected logger, if available.
	Log(format string, v ...any)
}
