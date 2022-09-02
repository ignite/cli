package cosmosgen

import (
	"github.com/ignite-hq/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite-hq/cli/ignite/pkg/giturl"
	"github.com/ignite-hq/cli/ignite/pkg/gomodulepath"

	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type vuexGenerator struct {
	g *generator
}

type generateVuexPayload struct {
	Modules []module.Module
	User    string ``
	Repo    string ``
}

func newVuexGenerator(g *generator) *vuexGenerator {
	return &vuexGenerator{
		g: g,
	}
}

func (g *generator) generateVuex() error {
	vsg := newVuexGenerator(g)

	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	chainInfo, err := giturl.Parse(chainPath.RawPath)
	if err != nil {
		return err
	}

	data := generatePayload{
		Modules: g.appModules,
		User:    chainInfo.User,
		Repo:    chainInfo.Repo,
	}

	if g.o.jsIncludeThirdParty {
		for _, modules := range g.thirdModules {
			data.Modules = append(data.Modules, modules...)
		}
	}

	if err := vsg.generateVueTemplates(data); err != nil {
		return err
	}

	if err := vsg.generateRootTemplates(data); err != nil {
		return err
	}

	return nil
}

func (g *vuexGenerator) generateVueTemplates(payload generatePayload) error {
	gg := &errgroup.Group{}

	generate := func() {
		for _, m := range payload.Modules {
			m := m

			gg.Go(func() error {
				vueAPIOut := filepath.Join(g.g.o.vuexRootPath, m.Pkg.Name)

				if err := os.MkdirAll(vueAPIOut, 0766); err != nil {
					return err
				}

				if err := templateTSClientVue.Write(vueAPIOut, "", struct {
					Module module.Module
					User   string
					Repo   string
				}{
					Module: m,
					User:   payload.User,
					Repo:   payload.Repo,
				}); err != nil {
					return err
				}

				return nil
			})
		}
	}

	generate()

	return gg.Wait()
}

func (g *vuexGenerator) generateRootTemplates(payload generatePayload) error {
	vueOut := filepath.Join(g.g.o.vuexRootPath)
	if err := os.MkdirAll(vueOut, 0766); err != nil {
		return err
	}
	if err := templateTSClientVueRoot.Write(vueOut, "", payload); err != nil {
		return err
	}

	return nil
}
