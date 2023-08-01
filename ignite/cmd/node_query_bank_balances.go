package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
)

func NewNodeQueryBankBalances() *cobra.Command {
	c := &cobra.Command{
		Use:   "balances [from_account_or_address]",
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
	session := cliui.New(cliui.StartSpinnerWithText(statusQuerying))
	defer session.End()

	inputAccount := args[0]

	client, err := newNodeCosmosClient(cmd)
	if err != nil {
		return err
	}

	// inputAccount can be an account of the keyring or a raw address
	address, err := client.Address(inputAccount)
	if err != nil {
		address = inputAccount
	}

	pagination, err := getPagination(cmd)
	if err != nil {
		return err
	}

	balances, err := client.BankBalances(cmd.Context(), address, pagination)
	if err != nil {
		return err
	}

	var rows [][]string
	for _, b := range balances {
		rows = append(rows, []string{b.Amount.String(), b.Denom})
	}

	return session.PrintTable([]string{"Amount", "Denom"}, rows...)
}
