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

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	grouptypes "github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	appanalysis "github.com/ignite/cli/ignite/pkg/cosmosanalysis/app"
	"github.com/ignite/cli/ignite/pkg/cosmosver"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/validation"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field"
	"github.com/ignite/cli/ignite/templates/module"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
	moduleimport "github.com/ignite/cli/ignite/templates/module/import"
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
	// reservedNames are either names from the default modules defined in a Cosmos-SDK app or names used in the default query and tx CLI namespace.
	// A new module's name can't be equal to a reserved name.
	// A map is used for direct comparing.
	reservedNames = map[string]struct{}{
		"account":                    {},
		"block":                      {},
		"broadcast":                  {},
		"encode":                     {},
		"multisign":                  {},
		"sign":                       {},
		"tx":                         {},
		"txs":                        {},
		ibcexported.ModuleName:       {},
		transfertypes.ModuleName:     {},
		authtypes.ModuleName:         {},
		authztypes.ModuleName:        {},
		banktypes.ModuleName:         {},
		crisistypes.ModuleName:       {},
		capabilitytypes.ModuleName:   {},
		distributiontypes.ModuleName: {},
		evidencetypes.ModuleName:     {},
		feegranttypes.ModuleName:     {},
		genutiltypes.ModuleName:      {},
		govtypes.ModuleName:          {},
		grouptypes.ModuleName:        {},
		minttypes.ModuleName:         {},
		paramstypes.ModuleName:       {},
		slashingtypes.ModuleName:     {},
		stakingtypes.ModuleName:      {},
		upgradetypes.ModuleName:      {},
		vestingtypes.ModuleName:      {},
	}

	// defaultStoreKeys are the names of the default store keys defined in a Cosmos-SDK app.
	// A new module's name can't have a defined store key in its prefix because of potential store key collision.
	defaultStoreKeys = []string{
		ibcexported.StoreKey,
		transfertypes.StoreKey,
		authtypes.StoreKey,
		banktypes.StoreKey,
		capabilitytypes.StoreKey,
		distributiontypes.StoreKey,
		evidencetypes.StoreKey,
		feegranttypes.StoreKey,
		govtypes.StoreKey,
		grouptypes.StoreKey,
		minttypes.StoreKey,
		paramstypes.StoreKey,
		slashingtypes.StoreKey,
		stakingtypes.StoreKey,
		upgradetypes.StoreKey,
	}
)

// moduleCreationOptions holds options for creating a new module.
type moduleCreationOptions struct {
	// ibc true if the module is an ibc module
	ibc bool

	// params list of parameters
	params []string

	// ibcChannelOrdering ibc channel ordering
	ibcChannelOrdering string

	// dependencies list of module dependencies
	dependencies []modulecreate.Dependency
}

// ModuleCreationOption configures Chain.
type ModuleCreationOption func(*moduleCreationOptions)

// WithIBC scaffolds a module with IBC enabled.
func WithIBC() ModuleCreationOption {
	return func(m *moduleCreationOptions) {
		m.ibc = true
	}
}

// WithParams scaffolds a module with params.
func WithParams(params []string) ModuleCreationOption {
	return func(m *moduleCreationOptions) {
		m.params = params
	}
}

// WithIBCChannelOrdering configures channel ordering of the IBC module.
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

// WithDependencies specifies the name of the modules that the module depends on.
func WithDependencies(dependencies []modulecreate.Dependency) ModuleCreationOption {
	return func(m *moduleCreationOptions) {
		m.dependencies = dependencies
	}
}

// CreateModule creates a new empty module in the scaffolded app.
func (s Scaffolder) CreateModule(
	ctx context.Context,
	cacheStorage cache.Storage,
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

	// Check if the module already exist
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return sm, err
	}
	if ok {
		return sm, fmt.Errorf("the module %v already exists", moduleName)
	}

	// Apply the options
	var creationOpts moduleCreationOptions
	for _, apply := range options {
		apply(&creationOpts)
	}

	// Parse params with the associated type
	params, err := field.ParseFields(creationOpts.params, checkForbiddenTypeIndex)
	if err != nil {
		return sm, err
	}

	// Check dependencies
	if err := checkDependencies(creationOpts.dependencies, s.path); err != nil {
		return sm, err
	}

	opts := &modulecreate.CreateOptions{
		ModuleName:   moduleName,
		ModulePath:   s.modpath.RawPath,
		Params:       params,
		AppName:      s.modpath.Package,
		AppPath:      s.path,
		IsIBC:        creationOpts.ibc,
		IBCOrdering:  creationOpts.ibcChannelOrdering,
		Dependencies: creationOpts.dependencies,
	}

	g, err := modulecreate.NewGenerator(opts)
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
	newSourceModification, runErr := xgenny.RunWithValidation(tracer, modulecreate.NewAppModify(tracer, opts))
	sm.Merge(newSourceModification)
	var validationErr validation.Error
	if runErr != nil && !errors.As(runErr, &validationErr) {
		return sm, runErr
	}

	return sm, finish(ctx, cacheStorage, opts.AppPath, s.modpath.RawPath)
}

// ImportModule imports specified module with name to the scaffolded app.
func (s Scaffolder) ImportModule(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	name string,
) (sm xgenny.SourceModification, err error) {
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
	g, err := moduleimport.NewGenerator(tracer, &moduleimport.ImportOptions{
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

	return sm, finish(ctx, cacheStorage, s.path, s.modpath.RawPath)
}

// moduleExists checks if the module exists in the app.
func moduleExists(appPath string, moduleName string) (bool, error) {
	absPath, err := filepath.Abs(filepath.Join(appPath, moduleDir, moduleName))
	if err != nil {
		return false, err
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		// The module doesn't exist
		return false, nil
	}

	return err == nil, err
}

// checkModuleName checks if the name can be used as a module name.
func checkModuleName(appPath, moduleName string) error {
	// go keyword
	if token.Lookup(moduleName).IsKeyword() {
		return fmt.Errorf("%s is a Go keyword", moduleName)
	}

	// check if the name is a reserved name
	if _, ok := reservedNames[moduleName]; ok {
		return fmt.Errorf("%s is a reserved name and can't be used as a module name", moduleName)
	}

	checkPrefix := func(name, prefix string) error {
		if strings.HasPrefix(name, prefix) {
			return fmt.Errorf("the module name can't be prefixed with %s because of potential store key collision", prefix)
		}
		return nil
	}

	// check if the name can imply potential store key collision
	for _, defaultStoreKey := range defaultStoreKeys {
		if err := checkPrefix(moduleName, defaultStoreKey); err != nil {
			return err
		}
	}

	// check store key with user's defined modules
	// we consider all user's defined modules use the module name as the store key
	entries, err := os.ReadDir(filepath.Join(appPath, moduleDir))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if err := checkPrefix(moduleName, entry.Name()); err != nil {
			return err
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
	switch {
	case s.Version.GTE(cosmosver.StargateFortyVersion):
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

// checkDependencies perform checks on the dependencies.
func checkDependencies(dependencies []modulecreate.Dependency, appPath string) error {
	depMap := make(map[string]struct{})
	for _, dep := range dependencies {
		// check the dependency has been registered
		path := filepath.Join(appPath, module.PathAppModule)
		if err := appanalysis.CheckKeeper(path, dep.KeeperName()); err != nil {
			return fmt.Errorf(
				"the module cannot have %s as a dependency: %w",
				dep.Name,
				err,
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
