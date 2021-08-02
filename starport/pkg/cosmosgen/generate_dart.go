package cosmosgen

import (
	"context"
	"os"

	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/protoc"
	protocgendart "github.com/tendermint/starport/starport/pkg/protoc-gen-dart"
	"golang.org/x/sync/errgroup"
)

var (
	dartOut = []string{
		"--dart_out=grpc:.",
	}
)

type dartGenerator struct {
	g *generator
}

func newDartGenerator(g *generator) *dartGenerator {
	return &dartGenerator{
		g: g,
	}
}

func (g *generator) generateDart() error {
	return newDartGenerator(g).generateModules()
}

func (g *dartGenerator) generateModules() error {
	flag, cleanup, err := protocgendart.Flag()
	if err != nil {
		return err
	}
	defer cleanup()

	gg := &errgroup.Group{}

	add := func(sourcePath string, modules []module.Module) {
		for _, m := range modules {
			m := m
			gg.Go(func() error { return g.generateModule(g.g.ctx, flag, sourcePath, m) })
		}
	}

	add(g.g.appPath, g.g.appModules)

	if g.g.o.dartIncludeThirdParty {
		for sourcePath, modules := range g.g.thirdModules {
			add(sourcePath, modules)
		}
	}

	return gg.Wait()
}

func (g *dartGenerator) generateModule(ctx context.Context, plugin, appPath string, m module.Module) error {
	out := g.g.o.dartOut(m)

	includePaths, err := g.g.resolveInclude(appPath)
	if err != nil {
		return err
	}

	// reset destination dir.
	if err := os.RemoveAll(out); err != nil {
		return err
	}
	if err := os.MkdirAll(out, 0766); err != nil {
		return err
	}

	return protoc.Generate(
		g.g.ctx,
		out,
		m.Pkg.Path,
		includePaths,
		dartOut,
		protoc.Plugin(plugin),
	)
}
