package cosmosver

import (
	"github.com/ignite/cli/ignite/pkg/gomodule"
)

const (
	// CosmosModulePath defines Cosmos SDK import path.
	CosmosModulePath = "github.com/cosmos/cosmos-sdk"
)

// Detect detects major version of Cosmos.
func Detect(appPath string) (version Version, err error) {
	parsed, err := gomodule.ParseAt(appPath)
	if err != nil {
		return version, err
	}

	for _, r := range parsed.Require {
		v := r.Mod

		if v.Path == CosmosModulePath {
			if version, err = Parse(v.Version); err != nil {
				return version, err
			}
		}
	}

	return
}
