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

	//go:embed genesistest/* genesistest/**/*
	fsGenesisTest embed.FS

	//go:embed simapp/* simapp/**/*
	fsSimapp embed.FS
)
