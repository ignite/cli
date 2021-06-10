package query

import (
	"github.com/tendermint/starport/starport/pkg/field"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

// Options ...
type Options struct {
	AppName     string
	ModuleName  string
	ModulePath  string
	OwnerName   string
	QueryName   multiformatname.MultiFormatName
	Description string
	ResFields   []field.Field
	ReqFields   []field.Field
	Paginated   bool
}
