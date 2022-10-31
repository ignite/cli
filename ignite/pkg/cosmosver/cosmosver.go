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

// Parse parses a Cosmos-SDK version.
func Parse(version string) (v Version, err error) {
	v.Version = version

	if v.Semantic, err = semver.Parse(strings.TrimPrefix(version, prefix)); err != nil {
		return v, err
	}

	return
}
