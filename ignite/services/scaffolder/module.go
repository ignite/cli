package scaffolder

import (
	"go/token"
	"os"
	"path/filepath"
	"strings"

	evidencetypes "cosmossdk.io/x/evidence/types"
	feegranttypes "cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	grouptypes "github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/gobuffalo/genny/v2"

	appanalysis "github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/app"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/validation"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field"
	"github.com/ignite/cli/v29/ignite/templates/module"
	modulecreate "github.com/ignite/cli/v29/ignite/templates/module/create"
)

const (
	moduleDir = "x"
	modulePkg = "module"
)

var (
	// reservedNames are either names from the default modules defined in a Cosmos-SDK app or names used in the default query and tx CLI namespace.
	// A new module's name can't be equal to a reserved name.
	// A map is used for direct comparing.
	reservedNames = map[string]struct{}{
		"account":                    {},
		"block":                      {},
		"broadcast":                  {},
		"consumer":                   {}, // ICS consumer module
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
	// ibc true if the module is an ibc module.
	ibc bool

	// params list of parameters.
	params []string

	// moduleConfigs list of module configs.
	moduleConfigs []string

	// ibcChannelOrdering ibc channel ordering.
	ibcChannelOrdering string

	// dependencies list of module dependencies.
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

// WithModuleConfigs scaffolds a module with module configs.
func WithModuleConfigs(moduleConfigs []string) ModuleCreationOption {
	return func(m *moduleCreationOptions) {
		m.moduleConfigs = moduleConfigs
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
	runner *xgenny.Runner,
	moduleName string,
	options ...ModuleCreationOption,
) error {
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return err
	}
	moduleName = mfName.LowerCase

	// Check if the module name is valid
	if err := checkModuleName(s.Path, moduleName); err != nil {
		return err
	}

	// Check if the module already exist
	ok, err := moduleExists(s.Path, moduleName)
	if err != nil {
		return err
	}
	if ok {
		return errors.Errorf("the module %v already exists", moduleName)
	}

	// Apply the options
	var creationOpts moduleCreationOptions
	for _, apply := range options {
		apply(&creationOpts)
	}

	// Parse params with the associated type
	params, err := field.ParseFields(creationOpts.params, checkForbiddenTypeIndex)
	if err != nil {
		return err
	}

	// Parse configs with the associated type
	configs, err := field.ParseFields(creationOpts.moduleConfigs, checkForbiddenTypeIndex)
	if err != nil {
		return err
	}

	// Check dependencies
	if err := checkDependencies(creationOpts.dependencies, s.Path); err != nil {
		return err
	}

	opts := &modulecreate.CreateOptions{
		ModuleName:   moduleName,
		ModulePath:   s.modpath.RawPath,
		Params:       params,
		Configs:      configs,
		AppName:      s.modpath.Package,
		AppPath:      s.Path,
		IsIBC:        creationOpts.ibc,
		IBCOrdering:  creationOpts.ibcChannelOrdering,
		Dependencies: creationOpts.dependencies,
	}

	g, err := modulecreate.NewGenerator(opts)
	if err != nil {
		return err
	}
	gens := []*genny.Generator{g}

	// Scaffold IBC module
	if opts.IsIBC {
		g, err = modulecreate.NewIBC(runner.Tracer(), opts)
		if err != nil {
			return err
		}
		gens = append(gens, g)
	}
	gens = append(gens, modulecreate.NewAppModify(runner.Tracer(), opts))

	err = runner.Run(gens...)
	var validationErr validation.Error
	if err != nil && !errors.As(err, &validationErr) {
		return err
	}
	return nil
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
		return errors.Errorf("%s is a Go keyword", moduleName)
	}

	// check if the name is a reserved name
	if _, ok := reservedNames[moduleName]; ok {
		return errors.Errorf("%s is a reserved name and can't be used as a module name", moduleName)
	}

	checkPrefix := func(name, prefix string) error {
		if strings.HasPrefix(name, prefix) {
			return errors.Errorf("the module name can't be prefixed with %s because of potential store key collision", prefix)
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

// checkDependencies perform checks on the dependencies.
func checkDependencies(dependencies []modulecreate.Dependency, appPath string) error {
	depMap := make(map[string]struct{})
	for _, dep := range dependencies {
		// check the dependency has been registered
		path := filepath.Join(appPath, module.PathAppModule)
		if err := appanalysis.CheckKeeper(path, dep.KeeperName()); err != nil {
			return errors.Errorf(
				"the module cannot have %s as a dependency: %w",
				dep.Name,
				err,
			)
		}

		// check duplicated
		_, ok := depMap[dep.Name]
		if ok {
			return errors.Errorf("%s is a duplicated dependency", dep)
		}
		depMap[dep.Name] = struct{}{}
	}

	return nil
}
