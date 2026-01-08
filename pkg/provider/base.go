package provider

type BaseProvider struct {
	logger Logger
}

func (b *BaseProvider) SetLogger(l Logger) {
	b.logger = l
}

func (b *BaseProvider) Log(format string, v ...any) {
	if b.logger != nil {
		b.logger.Printf(format, v...)
	}
}

func NewBaseProvider() *BaseProvider {
	return &BaseProvider{logger: &NoLogger{}}
}
