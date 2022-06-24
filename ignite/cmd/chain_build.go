package ignitecmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/services/chain"
)

const (
	flagOutput         = "output"
	flagRelease        = "release"
	flagReleaseTargets = "release.targets"
	flagReleasePrefix  = "release.prefix"
)

// NewChainBuild returns a new build command to build a blockchain app.
func NewChainBuild() *cobra.Command {
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
	- ignite chain build
	- ignite chain build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64`,
		Args: cobra.NoArgs,
		RunE: chainBuildHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetProto3rdParty("Available only without the --release flag"))
	c.Flags().Bool(flagRelease, false, "build for a release")
	c.Flags().StringSliceP(flagReleaseTargets, "t", []string{}, "release targets. Available only with --release flag")
	c.Flags().String(flagReleasePrefix, "", "tarball prefix for each release target. Available only with --release flag")
	c.Flags().StringP(flagOutput, "o", "", "binary output path")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")

	return c
}

func chainBuildHandler(cmd *cobra.Command, _ []string) error {
	var (
		isRelease, _      = cmd.Flags().GetBool(flagRelease)
		releaseTargets, _ = cmd.Flags().GetStringSlice(flagReleaseTargets)
		releasePrefix, _  = cmd.Flags().GetString(flagReleasePrefix)
		output, _         = cmd.Flags().GetString(flagOutput)
	)

	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	if flagGetProto3rdParty(cmd) {
		chainOption = append(chainOption, chain.EnableThirdPartyModuleCodegen())
	}

	c, err := newChainWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	if isRelease {
		releasePath, err := c.BuildRelease(cmd.Context(), cacheStorage, output, releasePrefix, releaseTargets...)
		if err != nil {
			return err
		}

		fmt.Printf("🗃  Release created: %s\n", colors.Info(releasePath))

		return nil
	}

	binaryName, err := c.Build(cmd.Context(), cacheStorage, output)
	if err != nil {
		return err
	}

	if output == "" {
		fmt.Printf("🗃  Installed. Use with: %s\n", colors.Info(binaryName))
	} else {
		binaryPath := filepath.Join(output, binaryName)
		fmt.Printf("🗃  Binary built at the path: %s\n", colors.Info(binaryPath))
	}

	return nil
}
