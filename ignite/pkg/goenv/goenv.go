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

	// GOPATH is the env var for GOPATH.
	GOPATH = "GOPATH"
)

const (
	binDir = "bin"
)

// Bin returns the path of where Go binaries are installed.
func Bin() string {
	if binPath := os.Getenv(GOBIN); binPath != "" {
		return binPath
	}
	if goPath := os.Getenv(GOPATH); goPath != "" {
		return filepath.Join(goPath, binDir)
	}
	return filepath.Join(build.Default.GOPATH, binDir)
}

// Path returns $PATH with correct go bin configuration set.
func Path() string {
	return os.ExpandEnv(fmt.Sprintf("$PATH:%s", Bin()))
}

// ConfigurePath configures the env with correct $PATH that has go bin setup.
func ConfigurePath() error {
	return os.Setenv("PATH", Path())
}
