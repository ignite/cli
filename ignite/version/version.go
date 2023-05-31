package version

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"text/tabwriter"

	"github.com/blang/semver/v4"
	"github.com/google/go-github/v48/github"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/cosmosver"
	"github.com/ignite/cli/ignite/pkg/gitpod"
	"github.com/ignite/cli/ignite/pkg/xexec"
)

const (
	errOldCosmosSDKVersion = `your chain has been scaffolded with an older version of Cosmos SDK: %s

Please, follow the migration guide to upgrade your chain to the latest version at https://docs.ignite.com/migration`

	versionDev     = "development"
	versionNightly = "nightly"
)

// Version is the semantic version of Ignite CLI.
var Version = versionDev

// CheckNext checks whether there is a new version of Ignite CLI.
func CheckNext(ctx context.Context) (isAvailable bool, version string, err error) {
	if Version == versionDev || Version == versionNightly {
		return false, "", nil
	}

	tagName, err := getLatestReleaseTag(ctx)
	if err != nil {
		return false, "", err
	}

	currentVersion, err := semver.ParseTolerant(Version)
	if err != nil {
		return false, "", err
	}

	latestVersion, err := semver.ParseTolerant(tagName)
	if err != nil {
		return false, "", err
	}

	isAvailable = latestVersion.GT(currentVersion)

	return isAvailable, tagName, nil
}

func getLatestReleaseTag(ctx context.Context) (string, error) {
	latest, _, err := github.
		NewClient(nil).
		Repositories.
		GetLatestRelease(ctx, "ignite", "cli")
	if err != nil {
		return "", err
	}

	if latest.TagName == nil {
		return "", nil
	}

	return *latest.TagName, nil
}

// resolveDevVersion creates a string for version printing if the version being used is "development".
// the version will be of the form "LATEST-dev" where LATEST is the latest tagged release.
func resolveDevVersion(ctx context.Context) string {
	// do nothing if built with specific tag
	if Version != versionDev && Version != versionNightly {
		return Version
	}

	tag, err := getLatestReleaseTag(ctx)
	if err != nil {
		return Version
	}

	if Version == versionDev {
		return tag + "-dev"
	}
	if Version == versionNightly {
		return tag + "-nightly"
	}

	return Version
}

// Long generates a detailed version info.
func Long(ctx context.Context) string {
	var (
		w          = &tabwriter.Writer{}
		b          = &bytes.Buffer{}
		date       = "undefined"
		head       = "undefined"
		modified   bool
		sdkVersion = "undefined"
	)
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == cosmosver.CosmosModulePath {
				sdkVersion = dep.Version
				break
			}
		}

		for _, kv := range info.Settings {
			switch kv.Key {
			case "vcs.revision":
				head = kv.Value
			case "vcs.time":
				date = kv.Value
			case "vcs.modified":
				modified = kv.Value == "true"
			}
		}
		if modified {
			// add * suffix to head to indicate the sources have been modified.
			head += "*"
		}
	}

	write := func(k string, v interface{}) {
		fmt.Fprintf(w, "%s:\t%v\n", k, v)
	}

	w.Init(b, 0, 8, 0, '\t', 0)

	write("Ignite CLI version", resolveDevVersion(ctx))
	write("Ignite CLI build date", date)
	write("Ignite CLI source hash", head)
	write("Ignite CLI config version", chainconfig.LatestVersion)
	write("Cosmos SDK version", sdkVersion)

	write("Your OS", runtime.GOOS)
	write("Your arch", runtime.GOARCH)

	cmdOut := &bytes.Buffer{}

	nodeJSCmd := "node"
	if xexec.IsCommandAvailable(nodeJSCmd) {
		cmdOut.Reset()

		err := exec.Exec(ctx, []string{nodeJSCmd, "-v"}, exec.StepOption(step.Stdout(cmdOut)))
		if err == nil {
			write("Your Node.js version", strings.TrimSpace(cmdOut.String()))
		}
	}

	cmdOut.Reset()
	err := exec.Exec(ctx, []string{"go", "version"}, exec.StepOption(step.Stdout(cmdOut)))
	if err != nil {
		panic(err)
	}
	write("Your go version", strings.TrimSpace(cmdOut.String()))

	unameCmd := "uname"
	if xexec.IsCommandAvailable(unameCmd) {
		cmdOut.Reset()

		err := exec.Exec(ctx, []string{unameCmd, "-a"}, exec.StepOption(step.Stdout(cmdOut)))
		if err == nil {
			write("Your uname -a", strings.TrimSpace(cmdOut.String()))
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		write("Your cwd", cwd)
	}

	write("Is on Gitpod", gitpod.IsOnGitpod())

	w.Flush()

	return b.String()
}

// AssertSupportedCosmosSDKVersion asserts that a Cosmos SDK version is supported by Ignite CLI.
func AssertSupportedCosmosSDKVersion(v cosmosver.Version) error {
	if v.LT(cosmosver.StargateFortySevenTwoVersion) {
		return fmt.Errorf(errOldCosmosSDKVersion, v)
	}
	return nil
}
