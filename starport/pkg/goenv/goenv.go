// Package goenv defines env variables known by Go and some utilities around it.
package goenv

import (
	"fmt"
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

// Path returns $PATH with correct go bin configuration set.
func Path() string {
	return os.ExpandEnv(fmt.Sprintf("$PATH:%s", GetGOBIN()))
}

// ConfigurePath configures the env with correct $PATH that has go bin setup.
func ConfigurePath() error {
	return os.Setenv("PATH", Path())
}
