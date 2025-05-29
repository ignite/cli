package modulecreate

import (
	"embed"
)

var (
	//go:embed files/base/* files/base/**/*
	fsBase embed.FS

	//go:embed files/ibc/* files/ibc/**/*
	fsIBC embed.FS
)
