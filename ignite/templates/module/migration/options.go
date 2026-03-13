package modulemigration

import (
	"fmt"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
)

// Options represents the options to scaffold a module migration.
type Options struct {
	ModuleName    string
	ModulePath    string
	MigrationName multiformatname.Name
	FromVersion   uint64
	ToVersion     uint64
}

// ModuleFile returns the path to the module definition file.
func (opts Options) ModuleFile() string {
	return filepath.Join("x", opts.ModuleName, "module", "module.go")
}

// MigrationVersion returns the migration package name.
func (opts Options) MigrationVersion() string {
	return fmt.Sprintf("v%d", opts.ToVersion)
}

// MigrationDir returns the path to the migration folder.
func (opts Options) MigrationDir() string {
	return filepath.Join("x", opts.ModuleName, "migrations", opts.MigrationVersion())
}

// MigrationFile returns the path to the migration source file.
func (opts Options) MigrationFile() string {
	return filepath.Join(opts.MigrationDir(), fmt.Sprintf("%s.go", opts.MigrationName.Snake))
}

// MigrationFunc returns the migration handler function name.
func (opts Options) MigrationFunc() string {
	return fmt.Sprintf("Migrate%s", opts.MigrationName.PascalCase)
}

// MigrationImportAlias returns the import alias used by module.go.
func (opts Options) MigrationImportAlias() string {
	return fmt.Sprintf("migrationv%d", opts.ToVersion)
}

// MigrationImportPath returns the migration import path used by module.go.
func (opts Options) MigrationImportPath() string {
	return fmt.Sprintf("%s/x/%s/migrations/%s", opts.ModulePath, opts.ModuleName, opts.MigrationVersion())
}
