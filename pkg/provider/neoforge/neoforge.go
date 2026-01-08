package neoforge

import "github.com/abulleDev/mcserverdl/pkg/provider"

type Provider struct {
	provider.BaseProvider
}

func New() *Provider {
	return &Provider{}
}
