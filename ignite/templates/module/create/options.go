package modulecreate

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/ignite/templates/field"
)

type (
	// ParamsOptions represents the options to scaffold a Cosmos SDK module parameters.
	ParamsOptions struct {
		ModuleName string
		AppName    string
		AppPath    string
		Params     field.Fields
	}

	// CreateOptions represents the options to scaffold a Cosmos SDK module.
	CreateOptions struct {
		ModuleName string
		ModulePath string
		AppName    string
		AppPath    string
		Params     field.Fields

		// True if the module should implement the IBC module interface
		IsIBC bool

		// Channel ordering of the IBC module: ordered, unordered or none
		IBCOrdering string

		// Dependencies of the module
		Dependencies Dependencies
	}

	// Dependency represents a module dependency of a module.
	Dependency struct {
		Name string
	}

	// Dependencies represents a list of module dependency.
	Dependencies []Dependency

	// MsgServerOptions defines options to add MsgServer.
	MsgServerOptions struct {
		ModuleName string
		ModulePath string
		AppName    string
		AppPath    string
	}
)

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

// Validate that options are usable.
func (opts *CreateOptions) Validate() error {
	return nil
}
