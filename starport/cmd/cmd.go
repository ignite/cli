package starportcmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/tendermint/starport/starport/internal/version"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/goenv"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

const flagHome = "home"
const checkVersionTimeout = time.Millisecond * 600

var (
	infoColor = color.New(color.FgYellow).SprintFunc()
)

// New creates a new root command for `starport` with its sub commands.
func New(ctx context.Context) *cobra.Command {
	cobra.EnableCommandSorting = false

	checkNewVersion(ctx)

	c := &cobra.Command{
		Use:   "starport",
		Short: "Starport offers everything you need to scaffold, test, build, and launch your blockchain",
		Long: `Starport is a tool for creating sovereign blockchains built with Cosmos SDK, the world’s
most popular modular blockchain framework. Starport offers everything you need to scaffold,
test, build, and launch your blockchain.

To get started, create a blockchain:

starport scaffold chain github.com/cosmonaut/mars`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return goenv.ConfigurePath()
		},
	}

	c.AddCommand(NewScaffold())
	c.AddCommand(NewChain())
	c.AddCommand(NewGenerate())
	c.AddCommand(NewNetwork())
	c.AddCommand(NewRelayer())
	c.AddCommand(NewTools())
	c.AddCommand(NewDocs())
	c.AddCommand(NewVersion())
	c.AddCommand(deprecated()...)

	return c
}

func logLevel(cmd *cobra.Command) chain.LogLvl {
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		return chain.LogVerbose
	}
	return chain.LogRegular
}

func printEvents(bus events.Bus, s *clispinner.Spinner) {
	for event := range bus {
		if event.IsOngoing() {
			s.SetText(event.Text())
			s.Start()
		} else {
			s.Stop()
			fmt.Printf("%s %s\n", color.New(color.FgGreen).SprintFunc()("✔"), event.Description)
		}
	}
}

func flagSetHome() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagHome, "", "Home directory used for blockchains")
	return fs
}

func getHomeFlag(cmd *cobra.Command) (home string) {
	home, _ = cmd.Flags().GetString(flagHome)
	return
}

func newChainWithHomeFlags(cmd *cobra.Command, appPath string, chainOption ...chain.Option) (*chain.Chain, error) {
	// Check if custom home is provided
	if home := getHomeFlag(cmd); home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}

	appPath, err := filepath.Abs(appPath)
	if err != nil {
		return nil, err
	}

	return chain.New(cmd.Context(), appPath, chainOption...)
}

func initOptionWithHomeFlag(cmd *cobra.Command, initOptions []networkbuilder.InitOption) []networkbuilder.InitOption {
	// Check if custom home is provided
	if home := getHomeFlag(cmd); home != "" {
		initOptions = append(initOptions, networkbuilder.InitializationHomePath(home))
	}

	return initOptions
}

var (
	modifyPrefix = color.New(color.FgMagenta).SprintFunc()("modify ")
	createPrefix = color.New(color.FgGreen).SprintFunc()("create ")
	removePrefix = func(s string) string {
		return strings.TrimPrefix(strings.TrimPrefix(s, modifyPrefix), createPrefix)
	}
)

func sourceModificationToString(sm xgenny.SourceModification) string {
	// get file names and add prefix
	var files []string
	for _, modified := range sm.ModifiedFiles() {
		files = append(files, modifyPrefix+modified)
	}
	for _, created := range sm.CreatedFiles() {
		files = append(files, createPrefix+created)
	}

	// sort filenames without prefix
	sort.Slice(files, func(i, j int) bool {
		s1 := removePrefix(files[i])
		s2 := removePrefix(files[j])

		return strings.Compare(s1, s2) == -1
	})

	return "\n" + strings.Join(files, "\n")
}

func deprecated() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:        "app",
			Deprecated: "use `starport scaffold chain` instead.",
		},
		{
			Use:        "build",
			Deprecated: "use `starport chain build` instead.",
		},
		{
			Use:        "serve",
			Deprecated: "use `starport chain serve` instead.",
		},
		{
			Use:        "faucet",
			Deprecated: "use `starport chain faucet` instead.",
		},
	}
}

func checkNewVersion(ctx context.Context) {
	if os.Getenv("GITPOD_WORKSPACE_ID") != "" {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, checkVersionTimeout)
	defer cancel()

	isAvailable, next, err := version.CheckNext(ctx)
	if err != nil || !isAvailable {
		return
	}

	fmt.Printf(`·
· 🛸 Starport %q is available!
·
· If you're looking to upgrade check out the instructions: https://docs.starport.network/intro/install.html#upgrading-your-starport-installation
·
··

`, next)
}
