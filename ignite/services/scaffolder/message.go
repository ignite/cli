package scaffolder

import (
	"context"
	"fmt"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field"
	"github.com/ignite/cli/ignite/templates/field/datatype"
	"github.com/ignite/cli/ignite/templates/message"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
)

// messageOptions represents configuration for the message scaffolding.
type messageOptions struct {
	description       string
	signer            string
	withoutSimulation bool
}

// newMessageOptions returns a messageOptions with default options.
func newMessageOptions(messageName string) messageOptions {
	return messageOptions{
		description: fmt.Sprintf("Broadcast message %s", messageName),
		signer:      "creator",
	}
}

// MessageOption configures the message scaffolding.
type MessageOption func(*messageOptions)

// WithDescription provides a custom description for the message CLI command.
func WithDescription(desc string) MessageOption {
	return func(m *messageOptions) {
		m.description = desc
	}
}

// WithSigner provides a custom signer name for the message.
func WithSigner(signer string) MessageOption {
	return func(m *messageOptions) {
		m.signer = signer
	}
}

// WithoutSimulation disables generating messages simulation.
func WithoutSimulation() MessageOption {
	return func(m *messageOptions) {
		m.withoutSimulation = true
	}
}

// AddMessage adds a new message to scaffolded app.
func (s Scaffolder) AddMessage(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	moduleName,
	msgName string,
	fields,
	resFields []string,
	options ...MessageOption,
) (sm xgenny.SourceModification, err error) {
	// Create the options
	scaffoldingOpts := newMessageOptions(msgName)
	for _, apply := range options {
		apply(&scaffoldingOpts)
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = s.modpath.Package
	}
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.LowerCase

	name, err := multiformatname.NewName(msgName)
	if err != nil {
		return sm, err
	}

	if err := checkComponentValidity(s.path, moduleName, name, false); err != nil {
		return sm, err
	}

	// Check and parse provided fields
	if err := checkCustomTypes(ctx, s.path, s.modpath.Package, moduleName, fields); err != nil {
		return sm, err
	}
	parsedMsgFields, err := field.ParseFields(fields, checkForbiddenMessageField, scaffoldingOpts.signer)
	if err != nil {
		return sm, err
	}

	// Check and parse provided response fields
	if err := checkCustomTypes(ctx, s.path, s.modpath.Package, moduleName, resFields); err != nil {
		return sm, err
	}
	parsedResFields, err := field.ParseFields(resFields, checkGoReservedWord, scaffoldingOpts.signer)
	if err != nil {
		return sm, err
	}

	mfSigner, err := multiformatname.NewName(scaffoldingOpts.signer)
	if err != nil {
		return sm, err
	}

	var (
		g    *genny.Generator
		opts = &message.Options{
			AppName:      s.modpath.Package,
			AppPath:      s.path,
			ModulePath:   s.modpath.RawPath,
			ModuleName:   moduleName,
			MsgName:      name,
			Fields:       parsedMsgFields,
			ResFields:    parsedResFields,
			MsgDesc:      scaffoldingOpts.description,
			MsgSigner:    mfSigner,
			NoSimulation: scaffoldingOpts.withoutSimulation,
		}
	)

	// Check and support MsgServer convention
	var gens []*genny.Generator
	gens, err = supportMsgServer(
		gens,
		tracer,
		s.path,
		&modulecreate.MsgServerOptions{
			ModuleName: opts.ModuleName,
			ModulePath: opts.ModulePath,
			AppName:    opts.AppName,
			AppPath:    opts.AppPath,
		},
	)
	if err != nil {
		return sm, err
	}

	gens, err = supportSimulation(
		gens,
		opts.AppPath,
		opts.ModulePath,
		opts.ModuleName,
	)
	if err != nil {
		return sm, err
	}

	// Scaffold
	g, err = message.NewGenerator(tracer, opts)
	if err != nil {
		return sm, err
	}
	gens = append(gens, g)
	sm, err = xgenny.RunWithValidation(tracer, gens...)
	if err != nil {
		return sm, err
	}
	return sm, finish(ctx, cacheStorage, opts.AppPath, s.modpath.RawPath)
}

// checkForbiddenMessageField returns true if the name is forbidden as a message name.
func checkForbiddenMessageField(name string) error {
	mfName, err := multiformatname.NewName(name)
	if err != nil {
		return err
	}

	if mfName.LowerCase == datatype.TypeCustom {
		return fmt.Errorf("%s is used by the message scaffolder", name)
	}

	return checkGoReservedWord(name)
}
