package ignitecmd

import (
	"fmt"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/spf13/cobra"
)

func NewNodeQueryBankBalances() *cobra.Command {
	c := &cobra.Command{
		Use:   "balances [account-name|address]",
		Short: "Query for account balances by account name or address",
		RunE:  nodeQueryBankBalancesHandler,
		Args:  cobra.ExactArgs(1),
	}

	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetAccountPrefixes())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetPagination("all balances"))

	return c
}

func nodeQueryBankBalancesHandler(cmd *cobra.Command, args []string) error {
	inputAccount := args[0]

	session := cliui.New()
	defer session.Cleanup()

	client, err := newNodeCosmosClient(cmd)
	if err != nil {
		return err
	}

	address, err := client.Bech32Address(inputAccount)
	if err != nil {
		return err
	}

	pagination, err := getPagination(cmd)
	if err != nil {
		return err
	}

	session.StartSpinner("Querying...")
	balances, err := client.BankBalances(cmd.Context(), address, pagination)
	if err != nil {
		return err
	}

	var rows [][]string
	for _, b := range balances {
		rows = append(rows, []string{fmt.Sprintf("%s", b.Amount), b.Denom})
	}

	session.StopSpinner()
	return session.PrintTable([]string{"Amount", "Denom"}, rows...)
}
