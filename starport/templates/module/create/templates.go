package modulecreate

import (
	"embed"
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
)
