package typed

import (
	"github.com/tendermint/starport/starport/pkg/field"
)

// Options ...
type Options struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	TypeName   string
	Fields     []field.Field
	NoMessage  bool
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
