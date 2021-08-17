package modulecreate

import (
	"embed"

	"github.com/tendermint/starport/starport/pkg/xgenny"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS

	//go:embed ibc/* ibc/**/*
	fsIBC embed.FS

	//go:embed msgserver/* msgserver/**/*
	fsMsgServer embed.FS

	//go:embed genesis_test/* genesis_test/**/*
	fsGenesisTest embed.FS

	stargateTemplate  = xgenny.NewEmbedWalker(fsStargate, "stargate/")
	ibcTemplate       = xgenny.NewEmbedWalker(fsIBC, "ibc/")
	msgServerTemplate = xgenny.NewEmbedWalker(fsMsgServer, "msgserver/")
	genesisTestTemplate = xgenny.NewEmbedWalker(fsGenesisTest, "genesis_test/")
)
