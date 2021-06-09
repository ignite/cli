package query

import (
	"github.com/tendermint/starport/starport/pkg/field"
)

// Options ...
type Options struct {
	AppName     string
	ModuleName  string
	ModulePath  string
	OwnerName   string
	QueryName   string
	Description string
	ResFields   []field.Field
	ReqFields   []field.Field
	Paginated   bool
}
