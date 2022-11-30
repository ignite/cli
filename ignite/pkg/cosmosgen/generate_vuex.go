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

	chainconfig "github.com/ignite/cli/ignite/config/chain"
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
	// Init the path to the "vue" folders inside the app
	vuePath := filepath.Join(g.appPath, chainconfig.DefaultVuePath)
	packagesPath := filepath.Join(vuePath, "package.json")
	if _, err := os.Stat(packagesPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	// Read the Vue app package file
	vuePkgRaw, err := os.ReadFile(packagesPath)
	if err != nil {
		return err
	}

	var vuePkg map[string]interface{}

	if err := json.Unmarshal(vuePkgRaw, &vuePkg); err != nil {
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
	tsClientVueRelPath, err := filepath.Rel(vuePath, tsClientPath)
	if err != nil {
		return err
	}

	err = mergo.Merge(&vuePkg, map[string]interface{}{
		"dependencies": map[string]interface{}{
			tsClientName: fmt.Sprintf("file:%s", tsClientVueRelPath),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to link ts-client dependency to the Vue app: %w", err)
	}

	// Save the modified package.json with the new dependencies
	vueFile, err := os.OpenFile(packagesPath, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer vueFile.Close()

	vueEnc := json.NewEncoder(vueFile)
	vueEnc.SetIndent("", "  ")
	vueEnc.SetEscapeHTML(false)
	if err := vueEnc.Encode(vuePkg); err != nil {
		return fmt.Errorf("error updating %s: %w", packagesPath, err)
	}

	return nil
}

func (g *generator) updateVuexDependencies() error {
	// Init the path to the "vuex" folders inside the app
	vuexPackagesPath := filepath.Join(g.o.vuexRootPath, "package.json")

	if _, err := os.Stat(vuexPackagesPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	// Read the Vuex stores package file
	vuexPkgRaw, err := os.ReadFile(vuexPackagesPath)
	if err != nil {
		return err
	}

	var vuexPkg map[string]interface{}

	if err := json.Unmarshal(vuexPkgRaw, &vuexPkg); err != nil {
		return fmt.Errorf("error parsing %s: %w", vuexPackagesPath, err)
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
	tsClientVuexRelPath, err := filepath.Rel(g.o.vuexRootPath, tsClientPath)
	if err != nil {
		return err
	}

	err = mergo.Merge(&vuexPkg, map[string]interface{}{
		"dependencies": map[string]interface{}{
			tsClientName: fmt.Sprintf("file:%s", tsClientVuexRelPath),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to link ts-client dependency to the Vuex stores: %w", err)
	}

	// Save the modified package.json with the new dependencies
	vuexFile, err := os.OpenFile(vuexPackagesPath, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer vuexFile.Close()

	vuexEnc := json.NewEncoder(vuexFile)
	vuexEnc.SetIndent("", "  ")
	vuexEnc.SetEscapeHTML(false)
	if err := vuexEnc.Encode(vuexPkg); err != nil {
		return fmt.Errorf("error updating %s: %w", vuexPackagesPath, err)
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

	for _, modules := range g.thirdModules {
		data.Modules = append(data.Modules, modules...)
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
