package scaffolder

import (
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	appanalysis "github.com/tendermint/starport/starport/pkg/cosmosanalysis/app"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gocmd"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/validation"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/module"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
	moduleimport "github.com/tendermint/starport/starport/templates/module/import"
)

const (
	wasmImport    = "github.com/CosmWasm/wasmd"
	wasmVersion   = "v0.16.0"
	extrasImport  = "github.com/tendermint/spm-extras"
	extrasVersion = "v0.1.0"
	appPkg        = "app"
	moduleDir     = "x"
)

var (
	// reservedNames are either names from the default modules defined in a Cosmos-SDK app or names used in the default query and tx CLI namespace
	// A new module's name can't be equal to a reserved name
	// A map is used for direct comparing
	reservedNames = map[string]struct{}{
		"account": {},
		"auth": {},
		"bank": {},
		"block": {},
		"broadcast": {},
		"crisis": {},
		"capability": {},
		"distribution": {},
		"encode": {},
		"evidence": {},
		"feegrant": {},
		"genutil": {},
		"gov": {},
		"group": {},
		"ibc": {},
		"mint": {},
		"multisign": {},
		"params": {},
		"sign": {},
		"slashing": {},
		"staking": {},
		"transfer": {},
		"tx": {},
		"txs": {},
		"upgrade": {},
		"vesting": {},
	}

	// defaultStoreKeys are the names of the default store keys defined in a Cosmos-SDK app
	// A new module's name can't have a defined store key in its prefix because of potential store key collision
	defaultStoreKeys = []string{
		"acc",
		"bank",
		"capability",
		"distribution",
		"evidence",
		"feegrant",
		"gov",
		"group",
		"mint",
		"slashing",
		"staking",
		"upgrade",
		"ibc",
		"transfer",
	}
)

// moduleCreationOptions holds options for creating a new module
type moduleCreationOptions struct {
	// chainID is the chain's id.
	ibc bool

	// homePath of the chain's config dir.
	ibcChannelOrdering string

	// list of module depencies
	dependencies []modulecreate.Dependency
}

// ModuleCreationOption configures Chain.
type ModuleCreationOption func(*moduleCreationOptions)

// WithIBC scaffolds a module with IBC enabled
func WithIBC() ModuleCreationOption {
	return func(m *moduleCreationOptions) {
		m.ibc = true
	}
}

// WithIBCChannelOrdering configures channel ordering of the IBC module
func WithIBCChannelOrdering(ordering string) ModuleCreationOption {
	return func(m *moduleCreationOptions) {
		switch ordering {
		case "ordered":
			m.ibcChannelOrdering = "ORDERED"
		case "unordered":
			m.ibcChannelOrdering = "UNORDERED"
		default:
			m.ibcChannelOrdering = "NONE"
		}
	}
}

// WithDependencies specifies the name of the modules that the module depends on
func WithDependencies(dependencies []modulecreate.Dependency) ModuleCreationOption {
	return func(m *moduleCreationOptions) {
		m.dependencies = dependencies
	}
}

// CreateModule creates a new empty module in the scaffolded app
func (s Scaffolder) CreateModule(
	tracer *placeholder.Tracer,
	moduleName string,
	options ...ModuleCreationOption,
) (sm xgenny.SourceModification, err error) {
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.LowerCase

	// Check if the module name is valid
	if err := checkModuleName(s.path, moduleName); err != nil {
		return sm, err
	}

	// Apply the options
	var creationOpts moduleCreationOptions
	for _, apply := range options {
		apply(&creationOpts)
	}

	// Check dependencies
	if err := checkDependencies(creationOpts.dependencies, s.path); err != nil {
		return sm, err
	}

	opts := &modulecreate.CreateOptions{
		ModuleName:   moduleName,
		ModulePath:   s.modpath.RawPath,
		AppName:      s.modpath.Package,
		AppPath:      s.path,
		OwnerName:    owner(s.modpath.RawPath),
		IsIBC:        creationOpts.ibc,
		IBCOrdering:  creationOpts.ibcChannelOrdering,
		Dependencies: creationOpts.dependencies,
	}

	// Generator from Cosmos SDK version
	g, err := modulecreate.NewStargate(opts)
	if err != nil {
		return sm, err
	}
	gens := []*genny.Generator{g}

	// Scaffold IBC module
	if opts.IsIBC {
		g, err = modulecreate.NewIBC(tracer, opts)
		if err != nil {
			return sm, err
		}
		gens = append(gens, g)
	}
	sm, err = xgenny.RunWithValidation(tracer, gens...)
	if err != nil {
		return sm, err
	}

	// Modify app.go to register the module
	newSourceModification, runErr := xgenny.RunWithValidation(tracer, modulecreate.NewStargateAppModify(tracer, opts))
	sm.Merge(newSourceModification)
	var validationErr validation.Error
	if runErr != nil && !errors.As(runErr, &validationErr) {
		return sm, runErr
	}

	return sm, finish(opts.AppPath, s.modpath.RawPath)
}

// ImportModule imports specified module with name to the scaffolded app.
func (s Scaffolder) ImportModule(tracer *placeholder.Tracer, name string) (sm xgenny.SourceModification, err error) {
	// Only wasm is currently supported
	if name != "wasm" {
		return sm, errors.New("module cannot be imported. Supported module: wasm")
	}

	ok, err := isWasmImported(s.path)
	if err != nil {
		return sm, err
	}
	if ok {
		return sm, errors.New("wasm is already imported")
	}

	// run generator
	g, err := moduleimport.NewStargate(tracer, &moduleimport.ImportOptions{
		AppPath:          s.path,
		Feature:          name,
		AppName:          s.modpath.Package,
		BinaryNamePrefix: s.modpath.Root,
	})
	if err != nil {
		return sm, err
	}

	sm, err = xgenny.RunWithValidation(tracer, g)
	if err != nil {
		var validationErr validation.Error
		if errors.As(err, &validationErr) {
			// TODO: implement a more generic method when there will be new methods to import wasm
			return sm, errors.New("wasm cannot be imported. Apps initialized with Starport <=0.16.2 must downgrade Starport to 0.16.2 to import wasm")
		}
		return sm, err
	}

	// import a specific version of ComsWasm
	// NOTE(dshulyak) it must be installed after validation
	if err := s.installWasm(); err != nil {
		return sm, err
	}

	return sm, finish(s.path, s.modpath.RawPath)
}

func checkModuleName(appPath, moduleName string) error {
	// go keyword
	if token.Lookup(moduleName).IsKeyword() {
		return fmt.Errorf("%s is a Go keyword", moduleName)
	}

	// check if the name is a reserved name
	if _, ok := reservedNames[moduleName]; ok {
		return fmt.Errorf("%s is a reserved name and can't be used as a module name", moduleName)
	}

	// check if the name can imply potential store key collision
	for _, defaultStoreKey := range defaultStoreKeys {
		if strings.HasPrefix(moduleName, defaultStoreKey) {
			return fmt.Errorf("the module name can't be prefixed with %s because of potential store key collision", defaultStoreKey)
		}
	}

	// check store key with user's defined modules
	// we consider all user's defined modules use the module name as the store key
	entries, err := os.ReadDir(filepath.Join(appPath, moduleDir))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if strings.HasPrefix(moduleName, entry.Name()) {
			return fmt.Errorf("the module name can't be prefixed with %s because of potential store key collision", entry.Name())
		}
	}

	return nil
}

func isWasmImported(appPath string) (bool, error) {
	abspath := filepath.Join(appPath, appPkg)
	fset := token.NewFileSet()
	all, err := parser.ParseDir(fset, abspath, func(os.FileInfo) bool { return true }, parser.ImportsOnly)
	if err != nil {
		return false, err
	}
	for _, pkg := range all {
		for _, f := range pkg.Files {
			for _, imp := range f.Imports {
				if strings.Contains(imp.Path.Value, wasmImport) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (s Scaffolder) installWasm() error {
	switch s.version {
	case cosmosver.StargateZeroFourtyAndAbove:
		return cmdrunner.
			New().
			Run(context.Background(),
				step.New(step.Exec(gocmd.Name(), "get", gocmd.PackageLiteral(wasmImport, wasmVersion))),
				step.New(step.Exec(gocmd.Name(), "get", gocmd.PackageLiteral(extrasImport, extrasVersion))),
			)
	default:
		return errors.New("version not supported")
	}
}

// checkDependencies perform checks on the dependencies
func checkDependencies(dependencies []modulecreate.Dependency, appPath string) error {
	depMap := make(map[string]struct{})
	for _, dep := range dependencies {
		// check the dependency has been registered
		path := filepath.Join(appPath, module.PathAppModule)
		if err := appanalysis.CheckKeeper(path, dep.KeeperName); err != nil {
			return fmt.Errorf(
				"the module cannot have %s as a dependency: %s",
				dep.Name,
				err.Error(),
			)
		}

		// check duplicated
		_, ok := depMap[dep.Name]
		if ok {
			return fmt.Errorf("%s is a duplicated dependency", dep)
		}
		depMap[dep.Name] = struct{}{}
	}

	return nil
}
