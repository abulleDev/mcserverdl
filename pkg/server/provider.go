package server

type Provider interface {
	GameVersions() ([]string, error)
	ServerVersions(gameVersion string) ([]string, error)
	DownloadURL(gameVersion, serverVersion string) (string, error)
}
