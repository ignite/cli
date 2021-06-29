package scaffolder

import (
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gocmd"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/validation"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/module"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
	moduleimport "github.com/tendermint/starport/starport/templates/module/import"
)

var (
	ErrNoIBCRouterPlaceholder = errors.New("app.go doesn't contain the necessary placeholder to generate an IBC module")
)

const (
	wasmImport  = "github.com/CosmWasm/wasmd"
	apppkg      = "app"
	moduleDir   = "x"
	wasmVersion = "v0.16.0"
)

// moduleCreationOptions holds options for creating a new module
type moduleCreationOptions struct {
	// chainID is the chain's id.
	ibc bool

	// homePath of the chain's config dir.
	ibcChannelOrdering string
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

// CreateModule creates a new empty module in the scaffolded app
func (s *Scaffolder) CreateModule(
	tracer *placeholder.Tracer,
	moduleName string,
	options ...ModuleCreationOption,
) (sm xgenny.SourceModification, err error) {
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.Lowercase

	// Check if the module name is valid
	if err := checkModuleName(moduleName); err != nil {
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
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return sm, err
	}
	opts := &modulecreate.CreateOptions{
		ModuleName:  moduleName,
		ModulePath:  path.RawPath,
		AppName:     path.Package,
		OwnerName:   owner(path.RawPath),
		IsIBC:       creationOpts.ibc,
		IBCOrdering: creationOpts.ibcChannelOrdering,
	}
	if opts.IsIBC {
		ibcPlaceholder, err := checkIBCRouterPlaceholder(s.path)
		if err != nil {
			return sm, err
		}
		if !ibcPlaceholder {
			return sm, ErrNoIBCRouterPlaceholder
		}
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

	newSourceModification, runErr := xgenny.RunWithValidation(tracer, modulecreate.NewStargateAppModify(tracer, opts))
	sm.Merge(newSourceModification)
	var validationErr validation.Error
	if runErr != nil && !errors.As(runErr, &validationErr) {
		return sm, runErr
	}

	// Generate proto and format the source
	pwd, err := os.Getwd()
	if err != nil {
		return sm, err
	}
	if err := s.finish(pwd, path.RawPath); err != nil {
		return sm, err
	}
	return sm, runErr
}

// ImportModule imports specified module with name to the scaffolded app.
func (s *Scaffolder) ImportModule(tracer *placeholder.Tracer, name string) (sm xgenny.SourceModification, err error) {
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

	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return sm, err
	}

	// run generator
	g, err := moduleimport.NewStargate(tracer, &moduleimport.ImportOptions{
		Feature:          name,
		AppName:          path.Package,
		BinaryNamePrefix: path.Root,
	})
	if err != nil {
		return sm, err
	}

	sm, err = xgenny.RunWithValidation(tracer, g)
	if err != nil {
		var validationErr validation.Error
		if errors.As(err, &validationErr) {
			// TODO: implement a more generic method when there will be new methods to import wasm
			return sm, errors.New("wasm cannot be imported. Apps initialized with Starport <=0.16.2 must downgrade Staport to 0.16.2 to import wasm")
		}
		return sm, err
	}

	// import a specific version of ComsWasm
	// NOTE(dshulyak) it must be installed after validation
	if err := s.installWasm(); err != nil {
		return sm, err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return sm, err
	}
	return sm, s.finish(pwd, path.RawPath)
}

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

	return true, err
}

func checkModuleName(moduleName string) error {
	// go keyword
	if token.Lookup(moduleName).IsKeyword() {
		return fmt.Errorf("%s is a Go keyword", moduleName)
	}

	// name of default registered module
	switch moduleName {
	case
		"auth",
		"genutil",
		"bank",
		"capability",
		"staking",
		"mint",
		"distr",
		"gov",
		"params",
		"crisis",
		"slashing",
		"ibc",
		"upgrade",
		"evidence",
		"transfer",
		"vesting":
		return fmt.Errorf("%s is a default module", moduleName)
	}
	return nil
}

func isWasmImported(appPath string) (bool, error) {
	abspath, err := filepath.Abs(filepath.Join(appPath, apppkg))
	if err != nil {
		return false, err
	}
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

func (s *Scaffolder) installWasm() error {
	switch s.version {
	case cosmosver.StargateZeroFourtyAndAbove:
		return cmdrunner.
			New(
				cmdrunner.DefaultStderr(os.Stderr),
			).
			Run(context.Background(),
				step.New(
					step.Exec(
						gocmd.Name(),
						"get",
						wasmImport+"@"+wasmVersion,
					),
				),
			)
	default:
		return errors.New("version not supported")
	}
}

// checkIBCRouterPlaceholder checks if app.go contains PlaceholderIBCAppRouter
// this placeholder is necessary to scaffold a new IBC module
// if it doesn't exist, we give instruction to add it to the user
func checkIBCRouterPlaceholder(appPath string) (bool, error) {
	appGo, err := filepath.Abs(filepath.Join(appPath, module.PathAppGo))
	if err != nil {
		return false, err
	}

	content, err := ioutil.ReadFile(appGo)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(content), module.PlaceholderIBCAppRouter), nil
}
