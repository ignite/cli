package cosmosgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
)

func (g *generator) checkVueExists() error {
	_, err := os.Stat(filepath.Join(g.appPath, g.frontendPath))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("frontend does not exist, please run `ignite scaffold vue` first")
	}

	return err
}

func (g *generator) updateComposableDependencies() error {
	if err := g.checkVueExists(); err != nil {
		return err
	}

	// Init the path to the appropriate frontend folder inside the app
	frontendPath := filepath.Join(g.appPath, g.frontendPath)
	packagesPath := filepath.Join(g.appPath, g.frontendPath, "package.json")

	b, err := os.ReadFile(packagesPath)
	if err != nil {
		return err
	}

	var pkg map[string]any
	if err := json.Unmarshal(b, &pkg); err != nil {
		return errors.Errorf("error parsing %s: %w", packagesPath, err)
	}

	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	// Make sure the TS client path is absolute
	tsClientPath, err := filepath.Abs(g.opts.tsClientRootPath)
	if err != nil {
		return errors.Errorf("failed to read the absolute typescript client path: %w", err)
	}

	// Add the link to the ts-client to the VUE app dependencies
	appModulePath := gomodulepath.ExtractAppPath(chainPath.RawPath)
	tsClientNS := strings.ReplaceAll(appModulePath, "/", "-")
	tsClientName := fmt.Sprintf("%s-client-ts", tsClientNS)
	tsClientRelPath, err := filepath.Rel(frontendPath, tsClientPath)
	if err != nil {
		return err
	}

	err = mergo.Merge(&pkg, map[string]interface{}{
		"dependencies": map[string]interface{}{
			tsClientName: fmt.Sprintf("file:%s", tsClientRelPath),
		},
	})
	if err != nil {
		return errors.Errorf("failed to link ts-client dependency in the frontend app: %w", err)
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
		return errors.Errorf("error updating %s: %w", packagesPath, err)
	}

	return nil
}

func (g *generator) generateComposables() error {
	if err := g.checkVueExists(); err != nil {
		return err
	}

	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	appModulePath := gomodulepath.ExtractAppPath(chainPath.RawPath)
	data := generatePayload{
		Modules:   g.appModules,
		PackageNS: strings.ReplaceAll(appModulePath, "/", "-"),
	}

	for _, modules := range g.thirdModules {
		data.Modules = append(data.Modules, modules...)
	}

	vsg := newComposablesGenerator(g)
	if err := vsg.generateComposableTemplates(data); err != nil {
		return err
	}

	return vsg.generateRootTemplates(data)
}

type composablesGenerator struct {
	g *generator
}

func newComposablesGenerator(g *generator) *composablesGenerator {
	return &composablesGenerator{g}
}

func (g *composablesGenerator) generateComposableTemplates(p generatePayload) error {
	gg := &errgroup.Group{}

	for _, m := range p.Modules {
		gg.Go(func() error {
			return g.generateComposableTemplate(m, p)
		})
	}

	return gg.Wait()
}

func (g *composablesGenerator) generateComposableTemplate(m module.Module, p generatePayload) error {
	outDir := g.g.opts.composablesOut(m)
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientComposable.Write(outDir, "", struct {
		Module    module.Module
		PackageNS string
	}{
		Module:    m,
		PackageNS: p.PackageNS,
	})
}

func (g *composablesGenerator) generateRootTemplates(p generatePayload) error {
	outDir := g.g.opts.composablesRootPath
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientComposableRoot.Write(outDir, "", p)
}
