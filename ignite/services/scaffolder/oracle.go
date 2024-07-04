package scaffolder

import (
	"context"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	"github.com/ignite/cli/v28/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/templates/ibc"
)

const (
	bandImport  = "github.com/bandprotocol/bandchain-packet"
	bandVersion = "v0.0.2"
)

// OracleOption configures options for AddOracle.
type OracleOption func(*oracleOptions)

type oracleOptions struct {
	signer string
}

// newOracleOptions returns a oracleOptions with default options
//
// Deprecated: This function is no longer maintained.
func newOracleOptions() oracleOptions {
	return oracleOptions{
		signer: "creator",
	}
}

// OracleWithSigner provides a custom signer name for the message
//
// Deprecated: This function is no longer maintained.
func OracleWithSigner(signer string) OracleOption {
	return func(m *oracleOptions) {
		m.signer = signer
	}
}

// AddOracle adds a new BandChain oracle envtest.
//
// Deprecated: This function is no longer maintained.
func (s *Scaffolder) AddOracle(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	moduleName,
	queryName string,
	options ...OracleOption,
) error {
	if err := s.installBandPacket(); err != nil {
		return err
	}

	o := newOracleOptions()
	for _, apply := range options {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return err
	}
	moduleName = mfName.LowerCase

	name, err := multiformatname.NewName(queryName)
	if err != nil {
		return err
	}

	if err := checkComponentValidity(s.path, moduleName, name, false); err != nil {
		return err
	}

	mfSigner, err := multiformatname.NewName(o.signer, checkForbiddenOracleFieldName)
	if err != nil {
		return err
	}

	// Module must implement IBC
	ok, err := isIBCModule(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("the module %s doesn't implement IBC module interface", moduleName)
	}

	// Generate the packet
	var (
		g    *genny.Generator
		opts = &ibc.OracleOptions{
			AppName:    s.modpath.Package,
			AppPath:    s.path,
			ModulePath: s.modpath.RawPath,
			ModuleName: moduleName,
			QueryName:  name,
			MsgSigner:  mfSigner,
		}
	)
	g, err = ibc.NewOracle(tracer, opts)
	if err != nil {
		return err
	}
	return s.runner.Run(g)
}

// Deprecated: This function is no longer maintained.
func (s Scaffolder) installBandPacket() error {
	return cmdrunner.New().
		Run(context.Background(),
			step.New(step.Exec(gocmd.Name(), "get", gocmd.PackageLiteral(bandImport, bandVersion))),
		)
}
