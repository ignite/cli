package scaffolder

import (
	"fmt"
	"github.com/tendermint/starport/starport/pkg/field"
	"os"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/message"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

// AddMessage adds a new message to scaffolded app
func (s *Scaffolder) AddMessage(
	tracer *placeholder.Tracer,
	moduleName,
	msgName,
	msgDesc string,
	fields,
	resFields []string,
) error {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = path.Package
	}

	if err := checkComponentValidity(s.path, moduleName, msgName); err != nil {
		return err
	}

	// Parse provided fields
	parsedMsgFields, err := field.ParseFields(fields, checkForbiddenMessageField)
	if err != nil {
		return err
	}
	parsedResFields, err := field.ParseFields(resFields, checkGoReservedWord)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &message.Options{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
			MsgName:    msgName,
			Fields:     parsedMsgFields,
			ResFields:  parsedResFields,
			MsgDesc:    msgDesc,
		}
	)

	// Check and support MsgServer convention
	var gens []*genny.Generator
	g, err = supportMsgServer(
		tracer,
		s.path,
		&modulecreate.MsgServerOptions{
			ModuleName: opts.ModuleName,
			ModulePath: opts.ModulePath,
			AppName:    opts.AppName,
			OwnerName:  opts.OwnerName,
		},
	)
	if err != nil {
		return err
	}
	if g != nil {
		gens = append(gens, g)
	}

	// Scaffold
	g, err = message.NewStargate(tracer, opts)
	if err != nil {
		return err
	}
	gens = append(gens, g)
	if err := xgenny.RunWithValidation(tracer, gens...); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return s.finish(pwd, path.RawPath)
}

// checkForbiddenMessageField returns true if the name is forbidden as a message name
func checkForbiddenMessageField(name string) error {
	if name == "creator" {
		return fmt.Errorf("%s is used by the message scaffolder", name)
	}

	return checkGoReservedWord(name)
}
