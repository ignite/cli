package cosmosver

import (
	"strings"

	"github.com/blang/semver/v4"
)

const prefix = "v"

// Version represents a range of Cosmos SDK versions.
type Version struct {
	// Version is the exact sdk version string.
	Version string

	// Semantic is the parsed version.
	Semantic semver.Version
}

var (
	StargateFortyVersion          = newVersion("0.40.0")
	StargateFortyFourVersion      = newVersion("0.44.0-alpha")
	StargateFortyFiveThreeVersion = newVersion("0.45.3")
	StargateFortySevenTwoVersion  = newVersion("0.47.2")
	StargateFiftyVersion          = newVersion("0.50.0") // 0.50.0 has been retracted and replaced with 0.50.1, but we keep it here for compatibility with pseudo versions.
	StargateFiftyThreeVersion     = newVersion("0.53.0")
)

var (
	// Versions is a list of known, sorted Cosmos-SDK versions.
	Versions = []Version{
		StargateFortyVersion,
		StargateFortyFourVersion,
		StargateFortyFiveThreeVersion,
		StargateFortySevenTwoVersion,
		StargateFiftyVersion,
		StargateFiftyThreeVersion, // NOTE: v0.50 and v0.53 are API compatible but consensus incompatible when using the new features from v0.53.
	}

	// Latest is the latest known version of the Cosmos-SDK.
	Latest = Versions[len(Versions)-1]
)

func newVersion(version string) Version {
	return Version{
		Version:  "v" + version,
		Semantic: semver.MustParse(version),
	}
}

// Parse parses a Cosmos-SDK version.
func Parse(version string) (v Version, err error) {
	v.Version = version

	if v.Semantic, err = semver.Parse(strings.TrimPrefix(version, prefix)); err != nil {
		return v, err
	}

	return
}

// GTE checks if v is greater than or equal to version.
func (v Version) GTE(version Version) bool {
	return v.Semantic.GTE(version.Semantic)
}

// LT checks if v is less than version.
func (v Version) LT(version Version) bool {
	return v.Semantic.LT(version.Semantic)
}

// LTE checks if v is less than or equal to version.
func (v Version) LTE(version Version) bool {
	return v.Semantic.LTE(version.Semantic)
}

// Is checks if v is equal to version.
func (v Version) Is(version Version) bool {
	return v.Semantic.EQ(version.Semantic)
}

func (v Version) String() string {
	return v.Version
}
