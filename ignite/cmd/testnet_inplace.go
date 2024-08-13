package ignitecmd

import (
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/spf13/cobra"
)

func NewTestNetInPlace() *cobra.Command {
	c := &cobra.Command{
		Use:   "in-place",
		Short: "Run simulation testing for the blockchain",
		Long:  "Run simulation testing for the blockchain. It sends many randomized-input messages of each module to a simulated node and checks if invariants break",
		Args:  cobra.NoArgs,
		RunE:  testnetInPlaceHandler,
	}
	flagSetPath(c)
	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	c.Flags().AddFlagSet(flagSetSkipProto())
	c.Flags().AddFlagSet(flagSetVerbose())

	c.Flags().Bool(flagQuitOnFail, false, "quit program if the app fails to start")
	return c
}

func testnetInPlaceHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.WithVerbosity(getVerbosity(cmd)),
	)
	defer session.End()

	// Otherwise run the serve command directly
	return chainInplace(cmd, session)
}

func chainInplace(cmd *cobra.Command, session *cliui.Session) error {
	chainOption := []chain.Option{
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

	cfg, err := c.Config()
	if err != nil {
		return err
	}

	var acc string
	for _, i := range cfg.Accounts {
		acc = acc + "," + i.Address
	}

	chainID, err := c.ID()
	if err != nil {
		return err
	}

	args := chain.InplaceArgs{
		NewChainID:         chainID,
		NewOperatorAddress: cfg.Validators[0].OperatorAddress,
		AcountsToFund:      acc,
	}
	return c.TestNetInPlace(cmd.Context(), args)
}
