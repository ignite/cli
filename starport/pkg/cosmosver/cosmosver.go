package cosmosver

import (
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/gomodule"
)

type MajorVersion string

const (
	// Launchpad points to Launchpad version of Cosmos-SDK.
	Launchpad MajorVersion = "launchpad"

	// Stargate points to Stargate version of Cosmos-SDK.
	Stargate MajorVersion = "stargate"
)

// MajorVersions are the list of supported Cosmos-SDK major versions.
var MajorVersions = majorVersions{Launchpad, Stargate}

const (
	tendermintPath                = "github.com/tendermint/tendermint"
	cosmosStargateTendermintMajor = "v0.34.0"
)

// Detect dedects major version of Cosmos.
func Detect(appPath string) (MajorVersion, error) {
	parsed, err := gomodule.ParseAt(appPath)
	if err != nil {
		return "", err
	}
	for _, r := range parsed.Require {
		v := r.Mod
		if v.Path == tendermintPath {
			if strings.HasPrefix(v.Version, cosmosStargateTendermintMajor) {
				return Stargate, nil
			}
			break
		}
	}
	return Launchpad, nil
}

type majorVersions []MajorVersion

// Parse checks if vs is a supported sdk version for scaffolding and if so,
// it parses it to sdkVersion.
func (v majorVersions) Parse(vs string) (MajorVersion, error) {
	for _, version := range v {
		if MajorVersion(vs) == version {
			return MajorVersion(vs), nil
		}
	}
	return "", fmt.Errorf("%q is an unkown sdk version", vs)
}

// String returns a string representation of the version list.
func (v majorVersions) String() string {
	var vs string
	for _, version := range v {
		vs += " -" + string(version)
	}
	return strings.TrimSpace(vs)
}
