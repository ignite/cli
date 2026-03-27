package scaffolder

import (
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	modulemigration "github.com/ignite/cli/v29/ignite/templates/module/migration"
)

// CreateModuleMigration scaffolds a new module migration inside an existing module.
func (s Scaffolder) CreateModuleMigration(moduleName string) error {
	mfModuleName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return err
	}
	moduleName = mfModuleName.LowerCase

	ok, err := moduleExists(s.appPath, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("the module %s doesn't exist", moduleName)
	}

	moduleFilePath := filepath.Join(s.appPath, moduleDir, moduleName, modulePkg, "module.go")
	content, err := os.ReadFile(moduleFilePath)
	if err != nil {
		return err
	}

	fromVersion, err := modulemigration.ConsensusVersion(string(content))
	if err != nil {
		return err
	}

	opts := &modulemigration.Options{
		ModuleName:  moduleName,
		ModulePath:  s.modpath.RawPath,
		FromVersion: fromVersion,
		ToVersion:   fromVersion + 1,
	}

	versionDir := filepath.Join(s.appPath, opts.MigrationDir())
	if _, err := os.Stat(versionDir); err == nil {
		return errors.Errorf("migration version %s already exists for module %s", opts.MigrationVersion(), moduleName)
	} else if !os.IsNotExist(err) {
		return err
	}

	g, err := modulemigration.NewGenerator(opts)
	if err != nil {
		return err
	}

	return s.Run(g)
}
