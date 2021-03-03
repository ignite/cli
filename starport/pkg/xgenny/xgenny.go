package xgenny

import (
	"bytes"
	"embed"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packd"
)

// Walker implements packd.Walker for Go embed's fs.FS.
type Walker struct {
	fs         embed.FS
	trimPrefix string
}

// NewEmbedWalker returns a new Walker for fs.
// trimPrefix is used to trim parent paths from the paths of found files.
func NewEmbedWalker(fs embed.FS, trimPrefix string) Walker {
	return Walker{fs, trimPrefix}
}

// Walk implements packd.Walker.
func (w Walker) Walk(wl packd.WalkFunc) error {
	return w.walkDir(wl, ".")
}

func (w Walker) walkDir(wl packd.WalkFunc, path string) error {
	entries, err := w.fs.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			w.walkDir(wl, filepath.Join(path, entry.Name()))
			continue
		}

		path := filepath.Join(path, entry.Name())

		data, err := w.fs.ReadFile(path)
		if err != nil {
			return err
		}

		ppath := strings.TrimPrefix(path, w.trimPrefix)

		f, err := packd.NewFile(ppath, bytes.NewReader(data))
		if err != nil {
			return err
		}

		wl(ppath, f)
	}

	return nil
}
