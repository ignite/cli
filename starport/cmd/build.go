package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/services/chain"
)

const (
	flagRebuildProtoOnce = "rebuild-proto-once"
	flagRelease          = "release"
	flagReleaseTargets   = "release.targets"
	flagReleasePrefix    = "release.prefix"
)

// NewBuild returns a new build command to build a blockchain app.
func NewBuild() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Build a node binary",
		Long: `By default, build your node binaries
and add the binaries to your $(go env GOPATH)/bin path.

To build binaries for a release, use the --release flag. The app binaries
for one or more specified release targets are built in a release/ dir under the app's
source. Specify the release targets with GOOS:GOARCH build tags.
If the optional --release.targets is not specified, a binary is created for your current environment.

Sample usages:
	- starport build
	- starport build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64`,
		Args: cobra.ExactArgs(0),
		RunE: buildHandler,
	}
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().StringVarP(&appPath, "path", "p", ".", "path of the app")
	c.Flags().Bool(flagRelease, false, "build for a release")
	c.Flags().StringSliceP(flagReleaseTargets, "t", []string{}, "release targets. Available only with --release flag")
	c.Flags().String(flagReleasePrefix, "", "tarball prefix for each release target. Available only with --release flag")
	c.Flags().Bool(flagRebuildProtoOnce, false, "Enables proto code generation for 3rd party modules. Available only without the --release flag")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func buildHandler(cmd *cobra.Command, args []string) error {
	isRebuildProtoOnce, err := cmd.Flags().GetBool(flagRebuildProtoOnce)
	if err != nil {
		return err
	}

	isRelease, _ := cmd.Flags().GetBool(flagRelease)
	releaseTargets, _ := cmd.Flags().GetStringSlice(flagReleaseTargets)
	releasePrefix, _ := cmd.Flags().GetString(flagReleasePrefix)

	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	if isRebuildProtoOnce {
		chainOption = append(chainOption, chain.EnableThirdPartyModuleCodegen())
	}

	c, err := newChainWithHomeFlags(cmd, appPath, chainOption...)
	if err != nil {
		return err
	}

	if isRelease {
		releasePath, err := c.BuildRelease(cmd.Context(), releasePrefix, releaseTargets...)
		if err != nil {
			return err
		}

		fmt.Printf("🗃  Release created: %s\n", infoColor(releasePath))

		return nil
	}

	binaryName, err := c.Build(cmd.Context())
	if err != nil {
		return err
	}

	fmt.Printf("🗃  Installed. Use with: %s\n", infoColor(binaryName))

	return nil
}
