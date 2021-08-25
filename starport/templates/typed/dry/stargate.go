package dry

import (
	"embed"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/typed"
)

var (
	//go:embed stargate/component/* stargate/component/**/*
	fsStargateComponent embed.FS

	stargateTemplate = xgenny.NewEmbedWalker(fsStargateComponent, "stargate/component/")
)

// NewStargate returns the generator to scaffold a basic type in a Stargate module.
func NewStargate(opts *typed.Options) (*genny.Generator, error) {
	g := genny.New()

	return g, typed.Box(stargateTemplate, opts, g)
}
