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

	"github.com/ignite/cli/v30/ignite/config"
	"github.com/ignite/cli/v30/ignite/pkg/errors"
)

// stdlibsFS embeds the gno standard library sources (from
// github.com/gnolang/gno) so ignite can boot a dev chain and run gno tests
// without requiring a local gno checkout.
//
// The tree is kept in a `_`-prefixed directory so the go tool ignores the
// copied .go files (they are gno data, not ignite code).
//
//go:embed all:_stdlibs/gnovm
var stdlibsFS embed.FS

// embeddedRoot is the embedded tree root within stdlibsFS. It holds the
// gnovm/stdlibs (chain stdlibs) and gnovm/tests/stdlibs (testing overrides).
const embeddedRoot = "_stdlibs/gnovm"

var (
	stdlibsMu      sync.Mutex // guards stdlib extraction across goroutines
	ensureHashV    string
	ensureHashOnce sync.Once
)

// stdlibsHash returns a fingerprint of the embedded stdlib tree.
func stdlibsHash() (string, error) {
	var err error
	ensureHashOnce.Do(func() {
		h := sha256.New()
		err = fs.WalkDir(stdlibsFS, embeddedRoot, func(path string, d fs.DirEntry, err error) error {
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
		if err == nil {
			ensureHashV = hex.EncodeToString(h.Sum(nil))
		}
	})
	if err != nil {
		return "", err
	}
	return ensureHashV, nil
}

// EnsureStdlibs extracts the embedded gno stdlibs under
// <ignite-config>/gno/gnovm/{stdlibs,tests/stdlibs} (if not already current)
// and returns the gno root directory for them, so the default
// <rootdir>/gnovm/... resolutions find both.
func EnsureStdlibs() (string, error) {
	stdlibsMu.Lock()
	defer stdlibsMu.Unlock()
	return extractStdlibs()
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
	gnovmDir := filepath.Join(gnoRoot, "gnovm")
	marker := stdlibMarkerPath(gnoRoot, hash)
	if _, err := os.Stat(marker); err == nil {
		return gnoRoot, nil // already extracted at this version
	}

	if err := writeStdlibTree(gnovmDir); err != nil {
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

// writeStdlibTree replaces gnovmDir with the embedded stdlib sources.
func writeStdlibTree(gnovmDir string) error {
	if err := os.RemoveAll(gnovmDir); err != nil {
		return err
	}
	return fs.WalkDir(stdlibsFS, embeddedRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// skip the embedded root itself; _stdlibs must not leak in
			if path == embeddedRoot {
				return nil
			}
			rel, err := filepath.Rel(embeddedRoot, path)
			if err != nil {
				return err
			}
			return os.MkdirAll(filepath.Join(gnovmDir, rel), 0o755)
		}
		rel, err := filepath.Rel(embeddedRoot, path)
		if err != nil {
			return err
		}
		target := filepath.Join(gnovmDir, rel)
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
