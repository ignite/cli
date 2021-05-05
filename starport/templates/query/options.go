package query

import "github.com/tendermint/starport/starport/templates/typed"

// Options ...
type Options struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	QueryName    string
	Description    string
	ResFields     []typed.Field
	ReqFields  []typed.Field
}

