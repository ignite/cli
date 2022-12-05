package cosmosgen

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
)

type composablesGenerator struct {
	g            *generator
	frontendType string
}

func newComposablesGenerator(g *generator, frontendType string) *composablesGenerator {
	return &composablesGenerator{g, frontendType}
}

func (g *generator) updateComposableDependencies(frontendType string) error {
	// Init the path to the appropriate frontend folder inside the app
	frontendPath := filepath.Join(g.appPath, frontendType)
	packagesPath := filepath.Join(frontendPath, "package.json")
	if _, err := os.Stat(packagesPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	// Read the Vue app package file
	b, err := os.ReadFile(packagesPath)
	if err != nil {
		return err
	}

	var pkg map[string]interface{}

	if err := json.Unmarshal(b, &pkg); err != nil {
		return fmt.Errorf("error parsing %s: %w", packagesPath, err)
	}

	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	// Make sure the TS client path is absolute
	tsClientPath, err := filepath.Abs(g.o.tsClientRootPath)
	if err != nil {
		return fmt.Errorf("failed to read the absolute typescript client path: %w", err)
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
		return fmt.Errorf("failed to link ts-client dependency in the frontend app: %w", err)
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

func (g *generator) generateComposables(frontendType string) error {
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

	vsg := newComposablesGenerator(g, frontendType)
	if err := vsg.generateComposableTemplates(data); err != nil {
		return err
	}

	return vsg.generateRootTemplates(data)
}

func (g *composablesGenerator) generateComposableTemplates(p generatePayload) error {
	gg := &errgroup.Group{}

	for _, m := range p.Modules {
		m := m

		gg.Go(func() error {
			return g.generateComposableTemplate(m, p)
		})
	}

	return gg.Wait()
}

func (g *composablesGenerator) generateComposableTemplate(m module.Module, p generatePayload) error {
	var outDir string
	if g.frontendType == "vue" {
		outDir = g.g.o.composablesOut(m)
	} else {
		outDir = g.g.o.hooksOut(m)
	}

	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientComposable.Write(outDir, "", struct {
		Module       module.Module
		PackageNS    string
		FrontendType string
	}{
		Module:       m,
		PackageNS:    p.PackageNS,
		FrontendType: g.frontendType,
	})
}

func (g *composablesGenerator) generateRootTemplates(p generatePayload) error {
	var outDir string
	if g.frontendType == "vue" {
		outDir = g.g.o.composablesRootPath
	} else {
		outDir = g.g.o.hooksRootPath
	}
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientComposableRoot.Write(outDir, "", p)
}
