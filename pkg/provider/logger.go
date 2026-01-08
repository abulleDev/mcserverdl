package provider

type Logger interface {
	Printf(format string, v ...any)
}

type NoLogger struct{}

func (*NoLogger) Printf(format string, v ...any) {}
