package ignitecmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/ignite/pkg/cliui"
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
	c.Flags().AddFlagSet(flagSetGenerateOnly())
	c.Flags().AddFlagSet(flagSetGasFlags())
	c.Flags().String(flagFees, "", "Fees to pay along with transaction; eg: 10uatom")

	return c
}

func nodeTxBankSendHandler(cmd *cobra.Command, args []string) error {
	var (
		fromAccountInput = args[0]
		toAccountInput   = args[1]
		amount           = args[2]
		generateOnly     = getGenerateOnly(cmd)
	)

	client, err := newNodeCosmosClient(cmd)
	if err != nil {
		return err
	}

	// fromAccountInput must be an account of the keyring
	fromAccount, err := client.Account(fromAccountInput)
	if err != nil {
		return err
	}

	// toAccountInput can be an account of the keyring or a raw address
	toAddress, err := client.Address(toAccountInput)
	if err != nil {
		toAddress = toAccountInput
	}

	coins, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return err
	}

	tx, err := client.BankSendTx(fromAccount, toAddress, coins)
	if err != nil {
		return err
	}

	session := cliui.New()
	defer session.Cleanup()
	if generateOnly {
		json, err := tx.EncodeJSON()
		if err != nil {
			return err
		}

		session.StopSpinner()
		return session.Println(string(json))
	}

	session.StartSpinner("Sending transaction...")
	if _, err := tx.Broadcast(); err != nil {
		return err
	}

	session.StopSpinner()
	session.Println("Transaction broadcast successful!")
	return session.Printf("%s sent from %s to %s\n", amount, fromAccountInput, toAccountInput)

}
