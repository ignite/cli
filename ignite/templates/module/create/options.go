package modulecreate

import (
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/v29/ignite/templates/field"
)

// ConfigsOptions represents the options to scaffold a Cosmos SDK module configs.
type ConfigsOptions struct {
	ModuleName string
	AppName    string
	ProtoDir   string
	ProtoVer   string
	Configs    field.Fields
}

// ProtoFile returns the path to the proto folder.
func (opts *ConfigsOptions) ProtoFile(fname string) string {
	return filepath.Join(opts.ProtoDir, opts.AppName, opts.ModuleName, opts.ProtoVer, fname)
}

// ParamsOptions represents the options to scaffold a Cosmos SDK module parameters.
type ParamsOptions struct {
	ModuleName string
	AppName    string
	ProtoDir   string
	ProtoVer   string
	Params     field.Fields
}

// ProtoFile returns the path to the proto folder.
func (opts *ParamsOptions) ProtoFile(fname string) string {
	return filepath.Join(opts.ProtoDir, opts.AppName, opts.ModuleName, opts.ProtoVer, fname)
}

// CreateOptions represents the options to scaffold a Cosmos SDK module.
type CreateOptions struct {
	ModuleName string
	ModulePath string
	AppName    string
	AppPath    string
	ProtoDir   string
	ProtoVer   string
	Params     field.Fields
	Configs    field.Fields

	// True if the module should implement the IBC module interface
	IsIBC bool

	// Channel ordering of the IBC module: ordered, unordered or none
	IBCOrdering string

	// Dependencies of the module
	Dependencies Dependencies
}

// ProtoFile returns the path to the proto folder.
func (opts *CreateOptions) ProtoFile(fname string) string {
	return filepath.Join(opts.ProtoDir, opts.AppName, opts.ModuleName, opts.ProtoVer, fname)
}

// Dependency represents a module dependency of a module.
type Dependency struct {
	Name string
}

// Dependencies represents a list of module dependency.
type Dependencies []Dependency

// NewDependency returns a new dependency.
func NewDependency(name string) Dependency {
	return Dependency{Name: strcase.ToCamel(name)}
}

// Contains returns true if contains dependency name.
func (d Dependencies) Contains(name string) bool {
	for _, dep := range d {
		if dep.Name == name {
			return true
		}
	}
	return false
}

// Len returns the length of dependencies.
func (d Dependencies) Len() int {
	return len(d)
}

// KeeperName returns the keeper's name for the dependency module.
func (d Dependency) KeeperName() string {
	return fmt.Sprint(d.Name, "Keeper")
}
