package doc

import (
	"embed"
	"fmt"
	"io/fs"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

//go:embed files/*
var files embed.FS

// Options represents the options to scaffold a migration document.
type Options struct {
	Path        string
	FromVersion *semver.Version
	ToVersion   *semver.Version
	Diffs       string
	Description string
}

func (o Options) position() string {
	return fmt.Sprintf("%02d%02d%02d", o.ToVersion.Major(), o.ToVersion.Minor(), o.ToVersion.Patch())
}

func (o Options) shortDescription() string {
	return fmt.Sprintf("Release %s", o.ToVersion.Original())
}

func (o Options) date() string {
	return time.Now().Format("Jan 2 15:04:05 2006")
}

// NewGenerator returns the generator to scaffold a migration doc.
func NewGenerator(opts Options) (*genny.Generator, error) {
	subFs, err := fs.Sub(files, "files")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	if err := g.OnlyFS(subFs, nil, nil); err != nil {
		return g, err
	}

	ctx := plush.NewContext()
	ctx.Set("Position", opts.position())
	ctx.Set("FromVersion", opts.FromVersion.Original())
	ctx.Set("ToVersion", opts.ToVersion.Original())
	ctx.Set("Diffs", opts.Diffs)
	ctx.Set("Description", opts.Description)
	ctx.Set("ShortDescription", opts.shortDescription())
	ctx.Set("Date", opts.date())

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{Version}}", opts.ToVersion.Original()))

	return g, nil
}
