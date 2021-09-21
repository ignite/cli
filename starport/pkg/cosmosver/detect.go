package cosmosver

import (
	"github.com/tendermint/starport/starport/pkg/gomodule"
)

const (
	cosmosModulePath = "github.com/cosmos/cosmos-sdk"
)

// Detect detects major version of Cosmos.
func Detect(appPath string) (version Version, err error) {
	parsed, err := gomodule.ParseAt(appPath)
	if err != nil {
		return version, err
	}
	for _, r := range parsed.Require {
		v := r.Mod
		if v.Path == cosmosModulePath {
			version, err = NewVersion(v.Version)
			if err != nil {
				return version, err
			}
		}
	}
	return
}
