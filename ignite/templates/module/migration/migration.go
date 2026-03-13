package modulemigration

import (
	"io/fs"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
)

// NewGenerator returns the generator to scaffold a new module migration.
func NewGenerator(opts *Options) (*genny.Generator, error) {
	subFS, err := fs.Sub(files, "files")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	if err := g.OnlyFS(subFS, nil, nil); err != nil {
		return g, err
	}

	ctx := plush.NewContext()
	ctx.Set("fromVersion", opts.FromVersion)
	ctx.Set("migrationFunc", opts.MigrationFunc())
	ctx.Set("migrationName", opts.MigrationName)
	ctx.Set("migrationVersion", opts.MigrationVersion())
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("toVersion", opts.ToVersion)

	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{migrationName}}", opts.MigrationName.Snake))
	g.Transformer(genny.Replace("{{migrationVersion}}", opts.MigrationVersion()))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.RunFn(moduleModify(opts))

	return g, nil
}
