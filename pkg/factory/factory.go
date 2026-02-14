package factory

import (
	"fmt"

	"github.com/abulleDev/mcserverdl/v2/pkg/provider"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/fabric"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/forge"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/neoforge"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/paper"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/purpur"
	"github.com/abulleDev/mcserverdl/v2/pkg/provider/vanilla"
)

func New(serverType string) (provider.Provider, error) {
	switch serverType {
	case "vanilla":
		return vanilla.New(), nil
	case "paper":
		return paper.New(), nil
	case "fabric":
		return fabric.New(), nil
	case "forge":
		return forge.New(), nil
	case "neoforge":
		return neoforge.New(), nil
	case "purpur":
		return purpur.New(), nil
	default:
		return nil, fmt.Errorf("unknown server type '%s'", serverType)
	}
}
