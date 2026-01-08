package vanilla

import "errors"

// ServerVersions returns the list of available server versions for a given game version.
// For vanilla, this always returns an error as there are no separate server versions.
func (p *Provider) ServerVersions(gameVersion string) ([]string, error) {
	p.Log("Vanilla does not support server versions")
	return nil, errors.New("vanilla server does not have version")
}
