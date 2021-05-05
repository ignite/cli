package query

import (
	"embed"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS

	stargateTemplate = xgenny.NewEmbedWalker(fsStargate, "stargate/")
)