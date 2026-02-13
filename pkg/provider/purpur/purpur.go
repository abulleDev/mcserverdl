package purpur

import "github.com/abulleDev/mcserverdl/v2/pkg/provider"

type Provider struct {
	provider.BaseProvider
}

func New() *Provider {
	return &Provider{}
}
