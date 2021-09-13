package cosmosver

import (
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"golang.org/x/mod/semver"
)

const (
	cosmosModulePath                  = "github.com/cosmos/cosmos-sdk"
	cosmosModuleMaxLaunchpadTag       = "v0.39.99"
	cosmosModuleStargateFortyTag      = "v0.40.0"
	cosmosModuleStargateFortyThreeTag = "v0.43.0"
)

// Detect detects major version of Cosmos.
func Detect(appPath string) (Version, error) {
	parsed, err := gomodule.ParseAt(appPath)
	if err != nil {
		return 0, err
	}
	for _, r := range parsed.Require {
		v := r.Mod
		if v.Path == cosmosModulePath {
			switch {
			case semver.Compare(v.Version, cosmosModuleStargateFortyThreeTag) >= 0:
				return StargateZeroFortyThreeAndAbove, nil

			case semver.Compare(v.Version, cosmosModuleStargateFortyTag) >= 0:
				return StargateZeroFortyToZeroFortyTwo, nil

			case semver.Compare(v.Version, cosmosModuleMaxLaunchpadTag) <= 0:
				return LaunchpadAny, nil

			default:
				return StargateBelowZeroForty, nil
			}
		}
	}
	return 0, nil
}
