package scaffolder

import (
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field"
	modulecreate "github.com/ignite/cli/v29/ignite/templates/module/create"
)

// CreateParams creates a new params in the scaffolded module.
func (s Scaffolder) CreateParams(
	moduleName string,
	params ...string,
) error {
	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = s.modpath.Package
	}
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return err
	}
	moduleName = mfName.LowerCase

	// Check if the module already exist
	ok, err := moduleExists(s.appPath, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("the module %v not exist", moduleName)
	}

	if err := checkParamCreated(s.appPath, moduleName, params); err != nil {
		return err
	}

	// Parse params with the associated type
	paramsFields, err := field.ParseFields(params, checkForbiddenTypeIndex)
	if err != nil {
		return err
	}

	opts := modulecreate.ParamsOptions{
		ModuleName: moduleName,
		Params:     paramsFields,
		AppName:    s.modpath.Package,
		ProtoDir:   s.protoDir,
		ProtoVer:   "v1", // TODO(@julienrbrt): possibly in the future add flag to specify custom proto version.
	}

	g, err := modulecreate.NewModuleParam(opts)
	if err != nil {
		return err
	}

	return s.Run(g)
}

// checkParamCreated checks if the parameter has been already created.
func checkParamCreated(appPath, moduleName string, params []string) error {
	path := filepath.Join(appPath, "x", moduleName, "types")
	ok, err := goanalysis.HasAnyStructFieldsInPkg(path, "Params", params)
	if err != nil {
		return err
	}

	if ok {
		return errors.Errorf(
			"duplicated params (%s) module %s",
			strings.Join(params, " "),
			moduleName,
		)
	}
	return nil
}
