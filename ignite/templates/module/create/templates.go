package modulecreate

import (
	"embed"
)

var (
	//go:embed files/base/* files/base/**/*
	fsBase embed.FS

	//go:embed files/ibc/* files/ibc/**/*
	fsIBC embed.FS

	//go:embed files/msgserver/* files/msgserver/**/*
	fsMsgServer embed.FS

	//go:embed files/genesistest/* files/genesistest/**/*
	fsGenesisTest embed.FS

	//go:embed files/simapp/* files/simapp/**/*
	fsSimapp embed.FS
)
