package ignitecmd

import (
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/goenv"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

const (
	flagCheckDependencies = "check-dependencies"
	flagDebug             = "debug"
	flagOutput            = "output"
	flagRelease           = "release"
	flagBuildTags         = "build.tags"
	flagReleasePrefix     = "release.prefix"
	flagReleaseTargets    = "release.targets"
)

// NewChainBuild returns a new build command to build a blockchain app.
func NewChainBuild() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Build a node binary",
		Long: `
The build command compiles the source code of the project into a binary and
installs the binary in the $(go env GOPATH)/bin directory.

You can customize the output directory for the binary using a flag:

	ignite chain build --output dist

To compile the binary Ignite first compiles protocol buffer (proto) files into
Go source code. Proto files contain required type and services definitions. If
you're using another program to compile proto files, you can use a flag to tell
Ignite to skip the proto compilation step:

	ignite chain build --skip-proto

Afterwards, Ignite install dependencies specified in the go.mod file. By default
Ignite doesn't check that dependencies of the main module stored in the module
cache have not been modified since they were downloaded. To enforce dependency
checking (essentially, running "go mod verify") use a flag:

	ignite chain build --check-dependencies

Next, Ignite identifies the "main" package of the project. By default the "main"
package is located in "cmd/{app}d" directory, where "{app}" is the name of the
scaffolded project and "d" stands for daemon. If your project contains more
than one "main" package, specify the path to the one that Ignite should compile
in config.yml:

	build:
	  main: custom/path/to/main

By default the binary name will match the top-level module name (specified in
go.mod) with a suffix "d". This can be customized in config.yml:

	build:
	  binary: mychaind

You can also specify custom linker flags:

	build:
	  ldflags:
	    - "-X main.Version=development"
	    - "-X main.Date=01/05/2022T19:54"

To build binaries for a release, use the --release flag. The binaries for one or
more specified release targets are built in a "release/" directory in the
project's source directory. Specify the release targets with GOOS:GOARCH build
tags. If the optional --release.targets is not specified, a binary is created
for your current environment.

	ignite chain build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64
`,
		Args: cobra.NoArgs,
		RunE: chainBuildHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	c.Flags().AddFlagSet(flagSetSkipProto())
	c.Flags().AddFlagSet(flagSetDebug())
	c.Flags().AddFlagSet(flagSetVerbose())
	c.Flags().Bool(flagRelease, false, "build for a release")
	c.Flags().StringSliceP(flagReleaseTargets, "t", []string{}, "release targets. Available only with --release flag")
	c.Flags().StringSlice(flagBuildTags, []string{}, "parameters to build the chain binary")
	c.Flags().String(flagReleasePrefix, "", "tarball prefix for each release target. Available only with --release flag")
	c.Flags().StringP(flagOutput, "o", "", "binary output path")

	return c
}

func chainBuildHandler(cmd *cobra.Command, _ []string) error {
	var (
		isRelease, _      = cmd.Flags().GetBool(flagRelease)
		releaseTargets, _ = cmd.Flags().GetStringSlice(flagReleaseTargets)
		releasePrefix, _  = cmd.Flags().GetString(flagReleasePrefix)
		buildTags, _      = cmd.Flags().GetStringSlice(flagBuildTags)
		output, _         = cmd.Flags().GetString(flagOutput)
		session           = cliui.New(
			cliui.WithVerbosity(getVerbosity(cmd)),
			cliui.StartSpinner(),
		)
	)
	defer session.End()

	chainOption := []chain.Option{
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
		chain.CheckCosmosSDKVersion(),
	}

	if flagGetCheckDependencies(cmd) {
		chainOption = append(chainOption, chain.CheckDependencies())
	}

	// check if custom config is defined
	config, _ := cmd.Flags().GetString(flagConfig)
	if config != "" {
		chainOption = append(chainOption, chain.ConfigFile(config))
	}

	c, err := chain.NewWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	if isRelease {
		releasePath, err := c.BuildRelease(ctx, cacheStorage, buildTags, output, releasePrefix, releaseTargets...)
		if err != nil {
			return err
		}

		return session.Printf("🗃  Release created: %s\n", colors.Info(releasePath))
	}

	binaryName, err := c.Build(ctx, cacheStorage, buildTags, output, flagGetSkipProto(cmd), flagGetDebug(cmd))
	if err != nil {
		return err
	}

	if output == "" {
		session.Printf("🗃  Installed. Use with: %s\n", colors.Info(binaryName))

		if _, err := exec.LookPath(binaryName); err != nil {
			session.Printf("⚠️  Warning: Binary not found in PATH\n")
			return session.Printf("   To run from anywhere, add Go bin to your PATH: export PATH=$PATH:%s\n", colors.Info(goenv.Bin()))
		}

		return nil
	}

	binaryPath := filepath.Join(output, binaryName)
	return session.Printf("🗃  Binary built at the path: %s\n", colors.Info(binaryPath))
}

func flagSetCheckDependencies() *flag.FlagSet {
	usage := "verify that cached dependencies have not been modified since they were downloaded"
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(flagCheckDependencies, false, usage)
	return fs
}

func flagGetCheckDependencies(cmd *cobra.Command) (check bool) {
	check, _ = cmd.Flags().GetBool(flagCheckDependencies)
	return
}

func flagSetDebug() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(flagDebug, false, "build a debug binary")
	return fs
}

func flagGetDebug(cmd *cobra.Command) (debug bool) {
	debug, _ = cmd.Flags().GetBool(flagDebug)
	return
}
