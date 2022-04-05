package query

import (
	"github.com/ignite-hq/cli/starport/pkg/multiformatname"
	"github.com/ignite-hq/cli/starport/templates/field"
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
