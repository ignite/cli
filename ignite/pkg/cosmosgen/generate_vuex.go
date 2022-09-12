package cosmosgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/imdario/mergo"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
)

type vuexGenerator struct {
	g *generator
}

func newVuexGenerator(g *generator) *vuexGenerator {
	return &vuexGenerator{g}
}

func (g *generator) updateVueDependencies() error {
	// Init the path to the "vue" folder inside the app
	vuePath := filepath.Join(g.appPath, "vue")
	packagesPath := filepath.Join(vuePath, "package.json")

	// Read the Vue app package file
	b, err := os.ReadFile(packagesPath)
	if err != nil {
		return err
	}

	var pkg map[string]interface{}

	if err := json.Unmarshal(b, &pkg); err != nil {
		return fmt.Errorf("error parsing %s: %w", packagesPath, err)
	}

	// Add the link to the ts-client to the VUE app dependencies
	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	appModulePath := gomodulepath.ExtractAppPath(chainPath.RawPath)
	tsClientNS := strings.ReplaceAll(appModulePath, "/", "-")
	tsClientName := fmt.Sprintf("%s-client-ts", tsClientNS)
	tsClientPath, err := filepath.Rel(vuePath, g.o.tsClientRootPath)
	if err != nil {
		return err
	}

	err = mergo.Merge(&pkg, map[string]interface{}{
		"dependencies": map[string]interface{}{
			tsClientName: fmt.Sprintf("file:%s", tsClientPath),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to link ts-client dependency in the Vue app: %w", err)
	}

	// Save the modified package.json with the new dependencies
	file, err := os.OpenFile(packagesPath, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(pkg); err != nil {
		return fmt.Errorf("error updating %s: %w", packagesPath, err)
	}

	return nil
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

	for _, m := range p.Modules {
		m := m

		gg.Go(func() error {
			return g.generateVueTemplate(m, p)
		})
	}

	return gg.Wait()
}

func (g *vuexGenerator) generateVueTemplate(m module.Module, p generatePayload) error {
	outDir := g.g.o.vuexOut(m)
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
