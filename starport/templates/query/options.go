package query

import (
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/templates/field"
)

// Options ...
type Options struct {
	AppName     string
	AppPath     string
	ModuleName  string
	ModulePath  string
	OwnerName   string
	QueryName   multiformatname.Name
	Description string
	ResFields   field.Fields
	ReqFields   field.Fields
	Paginated   bool
}
