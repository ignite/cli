package cosmosgen

import (
	"os"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
)

type vuexGenerator struct {
	g *generator
}

func newVuexGenerator(g *generator) *vuexGenerator {
	return &vuexGenerator{g}
}

func (g *generator) generateVuex() error {
	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	appModulePath := gomodulepath.ExtractAppPath(chainPath.RawPath)
	data := generatePayload{
		Modules:   g.appModules,
		PackageNS: strings.ReplaceAll(appModulePath, "/", "-"),
	}

	if g.o.jsIncludeThirdParty {
		for _, modules := range g.thirdModules {
			data.Modules = append(data.Modules, modules...)
		}
	}

	vsg := newVuexGenerator(g)
	if err := vsg.generateVueTemplates(data); err != nil {
		return err
	}

	return vsg.generateRootTemplates(data)
}

func (g *vuexGenerator) generateVueTemplates(p generatePayload) error {
	gg := &errgroup.Group{}

	func() {
		for _, m := range p.Modules {
			m := m

			gg.Go(func() error {
				return g.generateVueTemplate(m, p)
			})
		}
	}()

	return gg.Wait()
}

func (g *vuexGenerator) generateVueTemplate(m module.Module, p generatePayload) error {
	outDir := g.g.o.jsOut(m)
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientVue.Write(outDir, "", struct {
		Module    module.Module
		PackageNS string
	}{
		Module:    m,
		PackageNS: p.PackageNS,
	})
}

func (g *vuexGenerator) generateRootTemplates(p generatePayload) error {
	outDir := g.g.o.vuexRootPath
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientVueRoot.Write(outDir, "", p)
}
