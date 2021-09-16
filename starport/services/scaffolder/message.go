package scaffolder

import (
	"context"
	"fmt"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/message"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

// messageOptions represents configuration for the message scaffolding
type messageOptions struct {
	description string
	signer      string
}

// newMessageOptions returns a messageOptions with default options
func newMessageOptions(messageName string) messageOptions {
	return messageOptions{
		description: fmt.Sprintf("Broadcast message %s", messageName),
		signer:      "creator",
	}
}

// MessageOption configures the message scaffolding
type MessageOption func(*messageOptions)

// WithDescription provides a custom description for the message CLI command
func WithDescription(desc string) MessageOption {
	return func(m *messageOptions) {
		m.description = desc
	}
}

// WithSigner provides a custom signer name for the message
func WithSigner(signer string) MessageOption {
	return func(m *messageOptions) {
		m.signer = signer
	}
}

// AddMessage adds a new message to scaffolded app
func (s Scaffolder) AddMessage(
	ctx context.Context,
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
	if err := checkCustomTypes(ctx, s.path, moduleName, fields); err != nil {
		return sm, err
	}
	parsedMsgFields, err := field.ParseFields(fields, checkForbiddenMessageField)
	if err != nil {
		return sm, err
	}

	// Check and parse provided response fields
	if err := checkCustomTypes(ctx, s.path, moduleName, resFields); err != nil {
		return sm, err
	}
	parsedResFields, err := field.ParseFields(resFields, checkGoReservedWord)
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
			AppName:    s.modpath.Package,
			AppPath:    s.path,
			ModulePath: s.modpath.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(s.modpath.RawPath),
			MsgName:    name,
			Fields:     parsedMsgFields,
			ResFields:  parsedResFields,
			MsgDesc:    scaffoldingOpts.description,
			MsgSigner:  mfSigner,
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
			OwnerName:  opts.OwnerName,
		},
	)
	if err != nil {
		return sm, err
	}

	// Scaffold
	g, err = message.NewStargate(tracer, opts)
	if err != nil {
		return sm, err
	}
	gens = append(gens, g)
	sm, err = xgenny.RunWithValidation(tracer, gens...)
	if err != nil {
		return sm, err
	}
	return sm, finish(opts.AppPath, s.modpath.RawPath)
}

// checkForbiddenMessageField returns true if the name is forbidden as a message name
func checkForbiddenMessageField(name string) error {
	mfName, err := multiformatname.NewName(name)
	if err != nil {
		return err
	}

	switch mfName.LowerCase {
	case
		"creator",
		field.TypeCustom:
		return fmt.Errorf("%s is used by the packet scaffolder", name)
	}

	return checkGoReservedWord(name)
}
