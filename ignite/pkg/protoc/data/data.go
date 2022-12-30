package data

import (
	"embed"
	"io/fs"
)

//go:embed include/* include/**/*
var include embed.FS

// Include returns a file system that contains standard proto files used by protoc.
func Include() fs.FS {
	f, _ := fs.Sub(include, "include")
	return f
}

// Binary returns the platform specific protoc binary.
func Binary() []byte {
	return binary
}
