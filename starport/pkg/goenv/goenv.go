// Package goenv defines env variables known by Go and some utilities around it.
package goenv

import (
	"go/build"
	"os"
	"path/filepath"
)

const (
	// GOBIN is the env var for GOBIN.
	GOBIN = "GOBIN"
)

// GetGOBIN returns the path of where Go binaries are installed.
func GetGOBIN() string {
	if binPath := os.Getenv(GOBIN); binPath != "" {
		return binPath
	}
	return filepath.Join(build.Default.GOPATH, "bin")
}
