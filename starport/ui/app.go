package ui

import (
	"embed"
	"io/fs"
)

//go:embed app/dist/* app/dist/**/*
var ui embed.FS

// FS returns the file system containing the dev ui of Starport.
func FS() fs.FS {
	ufs, _ := fs.Sub(ui, "app/dist")
	return ufs
}
