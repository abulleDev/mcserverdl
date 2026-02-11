package vanilla

import (
	"context"
	"errors"
)

// ServerVersions returns the list of available server versions for a given game version.
// It uses a default background context.
// For vanilla, this always returns an error.
func (p *Provider) ServerVersions(gameVersion string) ([]string, error) {
	return p.ServerVersionsContext(context.Background(), gameVersion)
}

// ServerVersionsContext returns the list of available server versions for a given game version with context support.
// For vanilla, this always returns an error as there are no separate server versions.
//
// Parameters (unused for vanilla):
//   - ctx: the context to control the request lifetime.
//   - gameVersion: the Minecraft version string.
func (p *Provider) ServerVersionsContext(ctx context.Context, gameVersion string) ([]string, error) {
	p.Log("Vanilla does not support server versions")
	return nil, errors.New("vanilla server does not have version")
}
