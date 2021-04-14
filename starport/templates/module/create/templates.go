package modulecreate

import (
	"embed"

	"github.com/tendermint/starport/starport/pkg/xgenny"
)

// these needs to be created in the compiler time, otherwise packr2 won't be
// able to find boxes.
var (
	//go:embed launchpad/* launchpad/**/*
	fsLaunchpad embed.FS

	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS

	//go:embed ibc/* ibc/**/*
	fsIBC embed.FS

	//go:embed ibc/* ibc/**/*
	fsMsgServer embed.FS

	launchpadTemplate = xgenny.NewEmbedWalker(fsLaunchpad, "launchpad/")
	stargateTemplate  = xgenny.NewEmbedWalker(fsStargate, "stargate/")
	ibcTemplate       = xgenny.NewEmbedWalker(fsIBC, "ibc/")
	msgServerTemplate = xgenny.NewEmbedWalker(fsMsgServer, "msgserver/")
)
