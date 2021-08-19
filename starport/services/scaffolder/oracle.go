package scaffolder

import (
	"fmt"
	"os"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/ibc"
)

// OracleOption configures options for AddOracle.
type OracleOption func(*oracleOptions)

type oracleOptions struct {
	signer string
}

// newOracleOptions returns a oracleOptions with default options
func newOracleOptions() oracleOptions {
	return oracleOptions{
		signer: "creator",
	}
}

// OracleWithSigner provides a custom signer name for the message
func OracleWithSigner(signer string) OracleOption {
	return func(m *oracleOptions) {
		m.signer = signer
	}
}

// AddOracle adds a new BandChain oracle integration.
func (s *Scaffolder) AddOracle(
	tracer *placeholder.Tracer,
	moduleName,
	queryName string,
	options ...OracleOption,
) (sm xgenny.SourceModification, err error) {

	o := newOracleOptions()
	for _, apply := range options {
		apply(&o)
	}

	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return sm, err
	}

	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.Lowercase

	name, err := multiformatname.NewName(queryName)
	if err != nil {
		return sm, err
	}

	if err := checkComponentValidity(s.path, moduleName, name, false); err != nil {
		return sm, err
	}

	mfSigner, err := multiformatname.NewName(o.signer, checkForbiddenOracleFieldName)
	if err != nil {
		return sm, err
	}

	// Module must implement IBC
	ok, err := isIBCModule(s.path, moduleName)
	if err != nil {
		return sm, err
	}
	if !ok {
		return sm, fmt.Errorf("the module %s doesn't implement IBC module interface", moduleName)
	}

	// Generate the packet
	var (
		g    *genny.Generator
		opts = &ibc.OracleOptions{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
			QueryName:  name,
			MsgSigner:  mfSigner,
		}
	)
	g, err = ibc.NewOracle(tracer, opts)
	if err != nil {
		return sm, err
	}
	sm, err = xgenny.RunWithValidation(tracer, g)
	if err != nil {
		return sm, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return sm, err
	}
	return sm, s.finish(pwd, path.RawPath)
}
