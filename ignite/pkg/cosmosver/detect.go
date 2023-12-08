package cosmosver

import (
	"regexp"

	"github.com/ignite/cli/v28/ignite/pkg/gomodule"
)

var (
	// CosmosModulePath defines Cosmos SDK import path.
	CosmosModulePath = "cosmos-sdk"
	// CosmosSDK defines Cosmos SDK repository name
	CosmosSDK = "cosmos-sdk"
	// CosmosSDKModulePathPattern defines a regexp pattern for Cosmos SDK import path.
	CosmosSDKModulePathPattern = regexp.MustCompile(`github\.com\/[^\/]+\/cosmos-sdk`)
)

// Detect detects major version of Cosmos.
func Detect(appPath string) (version Version, err error) {
	parsed, err := gomodule.ParseAt(appPath)
	if err != nil {
		return version, err
	}

	for _, r := range parsed.Require {
		v := r.Mod

		if CosmosSDKModulePathPattern.MatchString(v.Path) {
			if version, err = Parse(v.Version); err != nil {
				return version, err
			}
		}
	}

	return
}
