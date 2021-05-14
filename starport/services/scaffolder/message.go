package scaffolder

import (
	"context"
	"fmt"
	"os"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/message"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

// AddMessage adds a new message to scaffolded app
func (s *Scaffolder) AddMessage(
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
	parsedMsgFields, err := parseFields(fields, checkForbiddenMessageField)
	if err != nil {
		return err
	}
	parsedResFields, err := parseFields(resFields, checkGoReservedWord)
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
	if err := supportMsgServer(
		s.path,
		&modulecreate.MsgServerOptions{
			ModuleName: opts.ModuleName,
			ModulePath: opts.ModulePath,
			AppName:    opts.AppName,
			OwnerName:  opts.OwnerName,
		},
	); err != nil {
		return err
	}

	// Scaffold
	g, err = message.NewStargate(opts)
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	if err := run.Run(); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := s.protoc(pwd, path.RawPath); err != nil {
		return err
	}
	return fmtProject(pwd)
}

// checkForbiddenMessageField returns true if the name is forbidden as a message name
func checkForbiddenMessageField(name string) error {
	if name == "creator" {
		return fmt.Errorf("%s is used by the message scaffolder", name)
	}

	return checkGoReservedWord(name)
}
