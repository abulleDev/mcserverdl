package provider

type Provider interface {
	GameVersions() ([]string, error)
	ServerVersions(gameVersion string) ([]string, error)
	DownloadURL(gameVersion, serverVersion string) (string, error)
	Download(gameVersion, serverVersion, installDir string, onProgress func(current, total int64)) error

	SetLogger(l Logger)
	Log(format string, v ...any)
}
