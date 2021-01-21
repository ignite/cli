package cosmosver

import (
	"fmt"
	"strings"
)

const (
	// Launchpad points to Launchpad version of Cosmos-SDK.
	Launchpad MajorVersion = "launchpad"

	// Stargate points to Stargate version of Cosmos-SDK.
	Stargate MajorVersion = "stargate"
)

const (
	LaunchpadAny Version = iota

	StargateBelowZeroFourty

	StargateZeroFourtyAndAbove
)

// MajorVersions are the list of supported Cosmos-SDK major versions.
var (
	MajorVersions = majorVersions{
		Launchpad,
		Stargate,
	}

	Versions = versions{
		LaunchpadAny,
		StargateBelowZeroFourty,
		StargateZeroFourtyAndAbove,
	}
)

// MajorVersion represents major, named versions of Cosmos-SDK.
type MajorVersion string

func (v MajorVersion) Is(comparedTo MajorVersion) bool {
	return v == comparedTo
}

// Version represents a range of Cosmos-SDK versions.
type Version int

func (v Version) Is(comparedTo Version) bool {
	return v == comparedTo
}

// Major returns the MajorVersion of the version.
func (v Version) Major() MajorVersion {
	switch v {
	case StargateBelowZeroFourty, StargateZeroFourtyAndAbove:
		return Stargate
	default:
		return Launchpad
	}
}

func (v Version) String() string {
	switch v {
	case StargateZeroFourtyAndAbove:
		return "Stargate v0.40.0 (or above)"

	case StargateBelowZeroFourty:
		return "Stargate v0.40.0 (pre-release)"

	default:
		return "Launchpad"
	}
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
	return "", fmt.Errorf("%q is an unknown sdk version", vs)
}

// String returns a string representation of the version list.
func (v majorVersions) String() string {
	var vs string
	for _, version := range v {
		vs += " -" + string(version)
	}
	return strings.TrimSpace(vs)
}

type versions []Version

func (v versions) Latest() Version {
	return v[len(v)-1]
}
