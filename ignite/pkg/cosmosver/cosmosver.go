package cosmosver

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
)

// Family represents the family(named versions) of Cosmos-SDK.
type Family string

const (
	// Launchpad represents the launchpad family of Cosmos-SDK.
	Launchpad Family = "launchpad"

	// Stargate represents the stargate family of Cosmos-SDK.
	Stargate Family = "stargate"
)

const prefix = "v"

// Version represents a range of Cosmos SDK versions.
type Version struct {
	// Family of the version
	Family Family

	// Version is the exact sdk version string.
	Version string

	// Semantic is the parsed version.
	Semantic semver.Version
}

var (
	MaxLaunchpadVersion           = newVersion("0.39.99", Launchpad)
	StargateFortyVersion          = newVersion("0.40.0", Stargate)
	StargateFortyFourVersion      = newVersion("0.44.0-alpha", Stargate)
	StargateFortyFiveThreeVersion = newVersion("0.45.3", Stargate)
)

var (
	// Versions is a list of known, sorted Cosmos-SDK versions.
	Versions = []Version{
		MaxLaunchpadVersion,
		StargateFortyVersion,
		StargateFortyFourVersion,
	}

	// Latest is the latest known version of the Cosmos-SDK.
	Latest = Versions[len(Versions)-1]
)

func newVersion(version string, family Family) Version {
	return Version{
		Family:   family,
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

	v.Family = Stargate
	if v.LTE(MaxLaunchpadVersion) {
		v.Family = Launchpad
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
	return fmt.Sprintf("%s - %s", v.Family, v.Version)
}

func (v Version) IsFamily(family Family) bool {
	return v.Family == family
}
