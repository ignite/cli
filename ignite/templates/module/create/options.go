package modulecreate

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/ignite/templates/field"
)

// CreateOptions represents the options to scaffold a Cosmos SDK module.
type CreateOptions struct {
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
	Dependencies []Dependency
}

// MsgServerOptions defines options to add MsgServer.
type MsgServerOptions struct {
	ModuleName string
	ModulePath string
	AppName    string
	AppPath    string
}

// Validate that options are usable.
func (opts *CreateOptions) Validate() error {
	return nil
}

// NewDependency returns a new dependency.
func NewDependency(name string) Dependency {
	return Dependency{Name: strcase.ToCamel(name)}
}

// Dependency represents a module dependency of a module.
type Dependency struct {
	Name string
}

// KeeperName returns the keeper's name for the dependency module.
func (d Dependency) KeeperName() string {
	return fmt.Sprint(d.Name, "Keeper")
}
