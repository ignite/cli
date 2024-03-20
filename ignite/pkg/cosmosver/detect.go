package cosmosver

import (
	"regexp"

	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
)

var (
	// CosmosSDKRepoName defines the name of the Cosmos SDK repository.
	CosmosSDKRepoName = "cosmos-sdk"
	// CosmosModulePath defines Cosmos SDK import path.
	CosmosModulePath = "github.com/cosmos/cosmos-sdk"
	// CosmosSDKModulePathPattern defines a regexp pattern for Cosmos SDK import path.
	CosmosSDKModulePathPattern = regexp.MustCompile(CosmosSDKRepoName + "$")
)

// Detect detects major version of Cosmos SDK.
// If the Cosmos SDK is replaced with a fork, it returns the version of the fork.
// If the Cosmos SDK is replaced with a local fork, it returns its non resolved version.
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
			// an empty version means that Cosmos SDK is replaced with a local fork
			// we fallback to use the non resolved go import of the Cosmos SDK
			if v.Version == "" {
				for _, r := range parsed.Require {
					if r.Mod.Path == CosmosModulePath {
						v.Version = r.Mod.Version
					}
				}
			}

			if version, err = Parse(v.Version); err != nil {
				return version, err
			}
		}
	}

	return
}
