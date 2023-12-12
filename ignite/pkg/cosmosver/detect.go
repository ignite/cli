package cosmosver

import (
	"regexp"

	"github.com/ignite/cli/v28/ignite/pkg/gomodule"
)

var (
	// CosmosModulePath defines Cosmos SDK import path.
	CosmosModulePath = "github.com/cosmos/cosmos-sdk"
	// CosmosSDKModulePathPattern defines a regexp pattern for Cosmos SDK import path.
	CosmosSDKModulePathPattern = regexp.MustCompile(`github\.com\/[^\/]+\/cosmos-sdk`)
)

// Detect detects major version of Cosmos SDK.
// If the Cosmos SDK is replaced with a fork, it will return the version of the fork.
func Detect(appPath string) (version Version, err error) {
	parsed, err := gomodule.ParseAt(appPath)
	if err != nil {
		return version, err
	}

	versions, err := gomodule.ResolveDependencies(parsed, false)
	if err != nil {
		return version, err
	}

	for _, v := range versions {
		if CosmosSDKModulePathPattern.MatchString(v.Path) {
			if version, err = Parse(v.Version); err != nil {
				return version, err
			}
		}
	}

	return
}
