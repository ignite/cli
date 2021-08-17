package typed

import (
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

// Options ...
type Options struct {
	AppName    string
	AppPath    string
	ModuleName string
	ModulePath string
	OwnerName  string
	TypeName   multiformatname.Name
	Fields     field.Fields
	NoMessage  bool
	Indexes    field.Fields
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
