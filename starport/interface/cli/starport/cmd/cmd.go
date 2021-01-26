package starportcmd

import (
	"fmt"
	"path/filepath"

	"github.com/tendermint/starport/starport/services/networkbuilder"

	flag "github.com/spf13/pflag"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/chain"
)

const (
	flagHome    = "home"
	flagCLIHome = "cli-home"
)

var (
	infoColor = color.New(color.FgYellow).SprintFunc()
)

// New creates a new root command for `starport` with its sub commands.
func New() *cobra.Command {
	c := &cobra.Command{
		Use:           "starport",
		Short:         "A tool for scaffolding out Cosmos applications",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	c.AddCommand(NewApp())
	c.AddCommand(NewType())
	c.AddCommand(NewServe())
	c.AddCommand(NewFaucet())
	c.AddCommand(NewBuild())
	c.AddCommand(NewModule())
	c.AddCommand(NewChain())
	c.AddCommand(NewRelayer())
	c.AddCommand(NewVersion())
	c.AddCommand(NewNetwork())
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

func flagSetHomes() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagHome, "", "Home directory used for blockchains")
	fs.String(flagCLIHome, "", "CLI home directory used for launchpad blockchains")
	return fs
}

func getHomeFlags(cmd *cobra.Command) (home string, cliHome string, err error) {
	home, err = cmd.Flags().GetString(flagHome)
	if err != nil {
		return home, cliHome, err
	}
	cliHome, err = cmd.Flags().GetString(flagCLIHome)
	return home, cliHome, err
}

func newChainWithHomeFlags(cmd *cobra.Command, appPath string, chainOption ...chain.Option) (*chain.Chain, error) {
	// Check if custom home is provided
	home, cliHome, err := getHomeFlags(cmd)
	if err != nil {
		return nil, err
	}
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	if cliHome != "" {
		chainOption = append(chainOption, chain.CLIHomePath(cliHome))
	}

	appPath, err = filepath.Abs(appPath)
	if err != nil {
		return nil, err
	}

	return chain.New(cmd.Context(), appPath, chainOption...)
}

func initOptionWithHomeFlags(cmd *cobra.Command, initOptions []networkbuilder.InitOption) ([]networkbuilder.InitOption, error) {
	// Check if custom home is provided
	home, cliHome, err := getHomeFlags(cmd)
	if err != nil {
		return initOptions, err
	}
	if home != "" {
		initOptions = append(initOptions, networkbuilder.InitializationHomePath(home))
	}
	if cliHome != "" {
		initOptions = append(initOptions, networkbuilder.InitializationCLIHomePath(cliHome))
	}

	return initOptions, nil
}
