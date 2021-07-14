package version

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/google/go-github/v37/github"
)

const versionDev = "development"
const prefix = "v"

var (
	// Version is the semantic version of Starport.
	Version = versionDev

	// Date is the build date of Starport.
	Date = ""

	// Head is the HEAD of the current branch.
	Head = ""
)

// CheckNext checks whether there is a new version of Starport.
func CheckNext(ctx context.Context) (isAvailable bool, version string, err error) {
	if Version == versionDev {
		return false, "", nil
	}

	latest, _, err := github.
		NewClient(nil).
		Repositories.
		GetLatestRelease(ctx, "tendermint", "starport")

	if err != nil {
		return false, "", err
	}

	if latest.TagName == nil {
		return false, "", nil
	}

	currentVersion, err := semver.Parse(strings.TrimPrefix(Version, prefix))
	if err != nil {
		return false, "", err
	}

	latestVersion, err := semver.Parse(strings.TrimPrefix(*latest.TagName, prefix))
	if err != nil {
		return false, "", err
	}

	isAvailable = latestVersion.GT(currentVersion)

	return isAvailable, *latest.TagName, nil
}

// Long generates a detailed version info.
func Long() string {
	output := fmt.Sprintf("starport version %s %s/%s -build date: %s",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		Date)

	if Head != "" {
		output = fmt.Sprintf("%s\ngit object hash: %s", output, Head)
	}
	return output
}
