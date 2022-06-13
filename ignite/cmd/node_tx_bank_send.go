package ignitecmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/spf13/cobra"
)

func NewNodeTxBankSend() *cobra.Command {
	c := &cobra.Command{
		Use:   "send [from_account_or_address] [to_account_or_address] [amount]",
		Short: "Send funds from one account to another.",
		RunE:  nodeTxBankSendHandler,
		Args:  cobra.ExactArgs(3),
	}

	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetAccountPrefixes())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetTxFrom())
	c.Flags().AddFlagSet(flagSetGenerateOnly())
	c.Flags().AddFlagSet(flagSetGasFlags())

	return c
}

func nodeTxBankSendHandler(cmd *cobra.Command, args []string) error {
	var (
		fromAccountInput = args[0]
		toAccountInput   = args[1]
		amount           = args[2]
		from             = getFrom(cmd)
		generateOnly     = getGenerateOnly(cmd)
	)

	session := cliui.New()
	defer session.Cleanup()

	session.StartSpinner("Sending transaction...")

	client, err := newNodeCosmosClient(cmd)
	if err != nil {
		return err
	}

	// If from flag is missing, check if the "from account" argument is an account name and use that instead
	if from == "" {
		fromInputIsAccount, err := client.AccountExists(fromAccountInput)
		if err != nil {
			return err
		}
		if fromInputIsAccount {
			from = fromAccountInput
		} else {
			return fmt.Errorf("\"--%s\" flag is required when from address is not an account name", flagFrom)
		}
	}

	fromAddress, err := client.Bech32Address(fromAccountInput)
	if err != nil {
		return err
	}

	toAddress, err := client.Bech32Address(toAccountInput)
	if err != nil {
		return err
	}

	coins, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return err
	}

	tx, err := client.BankSendTx(fromAddress, toAddress, coins, from)
	if err != nil {
		return err
	}

	if generateOnly {
		json, err := tx.EncodeJSON()
		if err != nil {
			return err
		}

		session.StopSpinner()
		return session.Println(string(json))
	}

	if _, err := tx.Broadcast(); err != nil {
		return err
	}

	session.StopSpinner()
	session.Println("Transaction broadcast successful!")
	return session.Printf("%s sent from %s to %s\n", amount, fromAccountInput, toAccountInput)

}
