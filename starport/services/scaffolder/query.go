package scaffolder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/message"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

// AddQuery adds a new query to scaffolded app
func (s *Scaffolder) AddQuery(moduleName string, queryName string, description string, resFields []string, reqFields []string) error {
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = path.Package
	}
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't exist", moduleName)
	}

	// Ensure the name is valid, otherwise it would generate an incorrect code
	if isForbiddenComponentName(queryName) {
		return fmt.Errorf("%s can't be used as a message name", queryName)
	}

	// Check component name is not already used
	ok, err = isComponentCreated(s.path, moduleName, queryName)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%s component is already added", queryName)
	}

	// Parse provided fields
	parsedReqFields, err := parseFields(reqFields, isGoReservedWord)
	if err != nil {
		return err
	}
	parsedResFields, err := parseFields(resFields, isGoReservedWord)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &message.Options{
			AppName:     path.Package,
			ModulePath:  path.RawPath,
			ModuleName:  moduleName,
			OwnerName:   owner(path.RawPath),
			QueryName:   queryName,
			ReqFields:   parsedReqFields,
			ResFields:   parsedResFields,
			Description: description,
		}
	)

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

// isQueryCreated checks if the message is already scaffolded TODO
func isQueryCreated(appPath, moduleName, msgName string) (isCreated bool, err error) {
	absPath, err := filepath.Abs(filepath.Join(
		appPath,
		moduleDir,
		moduleName,
		typesDirectory,
		"message_"+msgName+".go",
	))
	if err != nil {
		return false, err
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		// Message doesn't exist
		return false, nil
	}

	return true, err
}
