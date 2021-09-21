package cosmosver

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
)

type (
	// MajorVersion represents major, named versions of Cosmos-SDK.
	MajorVersion string

	// versions represents a list of Version
	versions []Version

	// Version represents a range of Cosmos SDK versions.
	Version struct {
		version  string
		Major    MajorVersion
		semantic semver.Version
	}
)

const (
	prefix = "v"

	// Launchpad points to Launchpad version of Cosmos-SDK.
	Launchpad MajorVersion = "launchpad"
	// Stargate points to Stargate version of Cosmos-SDK.
	Stargate MajorVersion = "stargate"
)

var (
	MaxLaunchpadVersion = Version{
		version:  "v0.39.99",
		semantic: semver.MustParse("0.39.99"),
		Major:    Launchpad,
	}
	StargateFortyVersion = Version{
		version:  "v0.40.0",
		semantic: semver.MustParse("0.40.0"),
		Major:    Stargate,
	}
	StargateFortyThreeVersion = Version{
		version:  "v0.43.0-alpha",
		semantic: semver.MustParse("0.43.0-alpha"),
		Major:    Stargate,
	}
	StargateFortyFourVersion = Version{
		version:  "v0.44.0",
		semantic: semver.MustParse("0.44.0"),
		Major:    Stargate,
	}

	// Versions are the list of supported Cosmos-SDK versions.
	Versions = versions{
		MaxLaunchpadVersion,
		StargateFortyVersion,
		StargateFortyThreeVersion,
		StargateFortyFourVersion,
	}
)

func NewVersion(version string) (v Version, err error) {
	v.version = version
	v.semantic, err = semver.Parse(strings.TrimPrefix(version, prefix))
	if err != nil {
		return v, err
	}

	v.Major = Stargate
	if v.LTE(MaxLaunchpadVersion) {
		v.Major = Launchpad
	}
	return
}

// GTE checks if v is greater than or equal to another.
func (v Version) GTE(version Version) bool {
	return v.semantic.GTE(version.semantic)
}

// LT checks if v is less than another.
func (v Version) LT(version Version) bool {
	return v.semantic.LT(version.semantic)
}

// LTE checks if v is less than or equal to another.
func (v Version) LTE(version Version) bool {
	return v.semantic.LTE(version.semantic)
}

// Is checks if v is equal to another.
func (v Version) Is(version Version) bool {
	return v.semantic.EQ(version.semantic)
}

func (v Version) String() string {
	return fmt.Sprintf("%s - %s", v.Major, v.version)
}

func (v Version) MajorIs(comparedTo MajorVersion) bool {
	return v.Major == comparedTo
}

func (v versions) Latest() Version {
	return v[len(v)-1]
}
