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

	//go:embed genesistest/module/* genesistest/module/**/*
	fsGenesisModuleTest embed.FS

	//go:embed genesistest/types/* genesistest/types/**/*
	fsGenesisTypesTest embed.FS

	stargateTemplate          = xgenny.NewEmbedWalker(fsStargate, "stargate/")
	ibcTemplate               = xgenny.NewEmbedWalker(fsIBC, "ibc/")
	msgServerTemplate         = xgenny.NewEmbedWalker(fsMsgServer, "msgserver/")
	genesisModuleTestTemplate = xgenny.NewEmbedWalker(fsGenesisModuleTest, "genesistest/module/")
	genesisTypesTestTemplate  = xgenny.NewEmbedWalker(fsGenesisTypesTest, "genesistest/types/")
)
