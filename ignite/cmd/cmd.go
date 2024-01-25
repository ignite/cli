package ignitecmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ignite/cli/v28/ignite/config"
	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	uilog "github.com/ignite/cli/v28/ignite/pkg/cliui/log"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gitpod"
	"github.com/ignite/cli/v28/ignite/pkg/goenv"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/version"
)

const (
	flagPath       = "path"
	flagHome       = "home"
	flagYes        = "yes"
	flagClearCache = "clear-cache"
	flagSkipProto  = "skip-proto"

	checkVersionTimeout = time.Millisecond * 600
	cacheFileName       = "ignite_cache.db"

	statusGenerating = "Generating..."
	statusQuerying   = "Querying..."
)

// New creates a new root command for `Ignite CLI` with its sub commands.
// Returns the cobra.Command, a cleanup function and an error. The cleanup
// function must be invoked by the caller to clean eventual Ignite App instances.
func New(ctx context.Context) (*cobra.Command, func(), error) {
	cobra.EnableCommandSorting = false

	c := &cobra.Command{
		Use:   "ignite",
		Short: "Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain",
		Long: `Ignite CLI is a tool for creating sovereign blockchains built with Cosmos SDK, the world’s
most popular modular blockchain framework. Ignite CLI offers everything you need to scaffold,
test, build, and launch your blockchain.

To get started, create a blockchain:

	ignite scaffold chain example
`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Check for new versions only when shell completion scripts are not being
			// generated to avoid invalid output to stdout when a new version is available
			if cmd.Use != "completion" {
				checkNewVersion(cmd.Context())
			}

			return goenv.ConfigurePath()
		},
	}

	c.AddCommand(
		NewScaffold(),
		NewChain(),
		NewGenerate(),
		NewNode(),
		NewAccount(),
		NewRelayer(),
		NewTools(),
		NewDocs(),
		NewVersion(),
		NewApp(),
		NewDoctor(),
		NewCompletionCmd(),
	)
	c.AddCommand(deprecated()...)

	// Load plugins if any
	if err := LoadPlugins(ctx, c); err != nil {
		return nil, nil, errors.Errorf("error while loading apps: %w", err)
	}
	return c, UnloadPlugins, nil
}

func getVerbosity(cmd *cobra.Command) uilog.Verbosity {
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		return uilog.VerbosityVerbose
	}

	return uilog.VerbosityDefault
}

func flagSetPath(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagPath, "p", ".", "path of the app")
}

func flagGetPath(cmd *cobra.Command) (path string) {
	path, _ = cmd.Flags().GetString(flagPath)
	return
}

func flagSetHome() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagHome, "", "directory where the blockchain node is initialized")
	return fs
}

func getHome(cmd *cobra.Command) (home string) {
	home, _ = cmd.Flags().GetString(flagHome)
	return
}

func flagSetConfig() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringP(flagConfig, "c", "", "path to Ignite config file (default: ./config.yml)")
	return fs
}

func getConfig(cmd *cobra.Command) (config string) {
	config, _ = cmd.Flags().GetString(flagConfig)
	return
}

func flagSetYes() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolP(flagYes, "y", false, "answers interactive yes/no questions with yes")
	return fs
}

func getYes(cmd *cobra.Command) (ok bool) {
	ok, _ = cmd.Flags().GetBool(flagYes)
	return
}

func flagSetSkipProto() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(flagSkipProto, false, "skip file generation from proto")
	return fs
}

func flagGetSkipProto(cmd *cobra.Command) bool {
	skip, _ := cmd.Flags().GetBool(flagSkipProto)
	return skip
}

func flagSetClearCache(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(flagClearCache, false, "clear the build cache (advanced)")
}

func flagGetClearCache(cmd *cobra.Command) bool {
	clearCache, _ := cmd.Flags().GetBool(flagClearCache)
	return clearCache
}

func newChainWithHomeFlags(cmd *cobra.Command, chainOption ...chain.Option) (*chain.Chain, error) {
	// Check if custom home is provided
	if home := getHome(cmd); home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}

	appPath := flagGetPath(cmd)
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return nil, err
	}

	return chain.New(absPath, chainOption...)
}

var (
	modifyPrefix = colors.Modified("modify ")
	createPrefix = colors.Success("create ")
	removePrefix = func(s string) string {
		return strings.TrimPrefix(strings.TrimPrefix(s, modifyPrefix), createPrefix)
	}
)

func sourceModificationToString(sm xgenny.SourceModification) (string, error) {
	// get file names and add prefix
	var files []string
	for _, modified := range sm.ModifiedFiles() {
		// get the relative app path from the current directory
		relativePath, err := relativePath(modified)
		if err != nil {
			return "", err
		}
		files = append(files, modifyPrefix+relativePath)
	}
	for _, created := range sm.CreatedFiles() {
		// get the relative app path from the current directory
		relativePath, err := relativePath(created)
		if err != nil {
			return "", err
		}
		files = append(files, createPrefix+relativePath)
	}

	// sort filenames without prefix
	sort.Slice(files, func(i, j int) bool {
		s1 := removePrefix(files[i])
		s2 := removePrefix(files[j])

		return strings.Compare(s1, s2) == -1
	})

	return "\n" + strings.Join(files, "\n"), nil
}

func deprecated() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:        "app",
			Deprecated: "use `ignite scaffold chain` instead.",
		},
		{
			Use:        "build",
			Deprecated: "use `ignite chain build` instead.",
		},
		{
			Use:        "serve",
			Deprecated: "use `ignite chain serve` instead.",
		},
		{
			Use:        "faucet",
			Deprecated: "use `ignite chain faucet` instead.",
		},
	}
}

// relativePath return the relative app path from the current directory.
func relativePath(appPath string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err := filepath.Rel(pwd, appPath)
	if err != nil {
		return "", err
	}
	return path, nil
}

func checkNewVersion(ctx context.Context) {
	if gitpod.IsOnGitpod() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, checkVersionTimeout)
	defer cancel()

	isAvailable, next, err := version.CheckNext(ctx)
	if err != nil || !isAvailable {
		return
	}

	fmt.Printf("⬆️ Ignite CLI %s is available! To upgrade: https://docs.ignite.com/welcome/install#upgrade", next)
}

func printSection(session *cliui.Session, title string) error {
	return session.Printf("------\n%s\n------\n\n", title)
}

func newCache(cmd *cobra.Command) (cache.Storage, error) {
	cacheRootDir, err := config.DirPath()
	if err != nil {
		return cache.Storage{}, err
	}

	storage, err := cache.NewStorage(
		filepath.Join(cacheRootDir, cacheFileName),
		cache.WithVersion(version.Version),
	)
	if err != nil {
		return cache.Storage{}, err
	}

	if flagGetClearCache(cmd) {
		if err := storage.Clear(); err != nil {
			return cache.Storage{}, err
		}
	}

	return storage, nil
}
