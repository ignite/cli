package gno

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ignite/cli/v29/ignite/config"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// stdlibsFS embeds the gno standard library sources (from
// github.com/gnolang/gno, gnovm/stdlibs) so ignite can boot a dev chain
// without requiring a local gno checkout.
//
// The tree is kept in a `_`-prefixed directory so the go tool ignores the
// copied .go files (they are gno data, not ignite code).
//
//go:embed all:_stdlibs/gnovm/stdlibs
var stdlibsFS embed.FS

// stdlibsRoot is the embedded tree root within stdlibsFS.
const stdlibsRoot = "_stdlibs/gnovm/stdlibs"

var (
	ensureOnce  sync.Once
	ensureErr   error
	ensureDir   string
	ensureHashV string
)

// stdlibsHash returns a fingerprint of the embedded stdlib tree.
func stdlibsHash() (string, error) {
	if ensureHashV != "" {
		return ensureHashV, nil
	}
	h := sha256.New()
	err := fs.WalkDir(stdlibsFS, stdlibsRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := stdlibsFS.ReadFile(path)
		if err != nil {
			return err
		}
		h.Write([]byte(path))
		h.Write(data)
		return nil
	})
	if err != nil {
		return "", err
	}
	ensureHashV = hex.EncodeToString(h.Sum(nil))
	return ensureHashV, nil
}

// EnsureStdlibs extracts the embedded gno stdlibs under
// <ignite-config>/gno/gnovm/stdlibs (if not already current) and returns the
// gno root directory to use as the node rootdir (so the default
// <rootdir>/gnovm/stdlibs resolution finds them).
func EnsureStdlibs() (string, error) {
	ensureOnce.Do(func() {
		ensureDir, ensureErr = extractStdlibs()
	})
	return ensureDir, ensureErr
}

// extractStdlibs writes the embedded stdlib tree to the ignite config dir
// when the current fingerprint is not already on disk.
func extractStdlibs() (string, error) {
	hash, err := stdlibsHash()
	if err != nil {
		return "", err
	}

	base, err := config.DirPath()
	if err != nil {
		return "", err
	}
	gnoRoot := filepath.Join(base, "gno")
	stdlibDir := filepath.Join(gnoRoot, "gnovm", "stdlibs")
	marker := stdlibMarkerPath(gnoRoot, hash)
	if _, err := os.Stat(marker); err == nil {
		return gnoRoot, nil // already extracted at this version
	}

	if err := writeStdlibTree(stdlibDir); err != nil {
		return "", errors.Errorf("extracting stdlibs: %w", err)
	}
	if err := os.WriteFile(marker, nil, 0o644); err != nil { //nolint:gosec // informational marker file
		return "", err
	}
	cleanStaleMarkers(gnoRoot, hash)
	return gnoRoot, nil
}

// stdlibMarkerPath returns the marker file recording an extracted version.
func stdlibMarkerPath(gnoRoot, hash string) string {
	return filepath.Join(gnoRoot, ".stdlibs-"+hash)
}

// writeStdlibTree replaces stdlibDir with the embedded stdlib sources.
func writeStdlibTree(stdlibDir string) error {
	if err := os.RemoveAll(stdlibDir); err != nil {
		return err
	}
	if err := os.MkdirAll(stdlibDir, 0o755); err != nil {
		return err
	}
	return fs.WalkDir(stdlibsFS, stdlibsRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(stdlibsRoot, path)
		if err != nil {
			return err
		}
		target := filepath.Join(stdlibDir, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		data, err := stdlibsFS.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644) //nolint:gosec // extracted stdlib sources are world-readable like the gno repo
	})
}

// cleanStaleMarkers removes marker files from older stdlib versions.
func cleanStaleMarkers(gnoRoot, currentHash string) {
	entries, err := os.ReadDir(gnoRoot)
	if err != nil {
		return
	}
	prefix := ".stdlibs-"
	current := prefix + currentHash
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, prefix) && name != current {
			_ = os.Remove(filepath.Join(gnoRoot, name))
		}
	}
}
