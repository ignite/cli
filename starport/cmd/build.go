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
		Short: "Build your app",
		Long: `By default, build command will build your application's binaries
and add them to your $(go env GOPATH)/bin path.

If you want to only build binaries for a release, use the --release flag with
--release.targets (optional). Then, binaries built for chosen targets will appear
in your app's path under release/. You can add any number of targets in GOOS:GOARCH format.
If you don't provide any, a binary will be created for your own machine.

Sample usages:
	- starport build
	- starport build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64`,
		Args: cobra.ExactArgs(0),
		RunE: buildHandler,
	}
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().StringVarP(&appPath, "path", "p", ".", "path of the app")
	c.Flags().Bool(flagRelease, false, "build for a release")
	c.Flags().StringSliceP(flagReleaseTargets, "t", []string{}, "release targets. only available with use of --release flag")
	c.Flags().String(flagReleasePrefix, "", "prefix to be used in tarball names for each target")
	c.Flags().Bool(flagRebuildProtoOnce, false, "Enables proto code generation for 3rd party modules. only available without the --release flag")
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

		fmt.Printf("🗃  Release created. Check: %s\n", infoColor(releasePath))

		return nil
	}

	binaryName, err := c.Build(cmd.Context())
	if err != nil {
		return err
	}

	fmt.Printf("🗃  Installed. Use with: %s\n", infoColor(binaryName))

	return nil
}
