package xgenny

import (
	"bytes"
	"embed"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/pkg/errors"

	"github.com/gobuffalo/packd"
)

// Walker implements packd.Walker for Go embed's fs.FS.
type Walker struct {
	fs         embed.FS
	trimPrefix string
	path       string
}

// NewEmbedWalker returns a new Walker for fs.
// trimPrefix is used to trim parent paths from the paths of found files.
func NewEmbedWalker(fs embed.FS, trimPrefix, path string) Walker {
	return Walker{fs: fs, trimPrefix: trimPrefix, path: path}
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

		entryPath := filepath.Join(path, entry.Name())

		data, err := w.fs.ReadFile(entryPath)
		if err != nil {
			return err
		}

		trimPath := strings.TrimPrefix(entryPath, w.trimPrefix)
		trimPath = filepath.Join(w.path, trimPath)
		f, err := packd.NewFile(trimPath, bytes.NewReader(data))
		if err != nil {
			return err
		}

		wl(trimPath, f)
	}

	return nil
}

// Transformer will plush-ify any file that has a ".plush" extension.
func Transformer(ctx *plush.Context) genny.Transformer {
	t := genny.NewTransformer(".plush", func(f genny.File) (genny.File, error) {
		s, err := plush.RenderR(f, ctx)
		if err != nil {
			return f, errors.Wrap(err, f.Name())
		}
		return genny.NewFileS(f.Name(), s), nil
	})
	t.StripExt = true
	return t
}
