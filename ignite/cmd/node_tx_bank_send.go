package ignitecmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
)

func NewNodeTxBankSend() *cobra.Command {
	c := &cobra.Command{
		Use:   "send [from_account_or_address] [to_account_or_address] [amount]",
		Short: "Send funds from one account to another.",
		RunE:  nodeTxBankSendHandler,
		Args:  cobra.ExactArgs(3),
	}

	return c
}

func nodeTxBankSendHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.End()

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

	tx, err := client.BankSendTx(cmd.Context(), fromAccount, toAddress, coins)
	if err != nil {
		return err
	}

	if generateOnly {
		json, err := tx.EncodeJSON()
		if err != nil {
			return err
		}

		return session.Println(string(json))
	}

	session.StartSpinner("Sending transaction...")
	resp, err := tx.Broadcast(cmd.Context())
	if err != nil {
		return err
	}

	session.Printf("Transaction broadcast successful! (hash = %s)\n", resp.TxHash)
	session.Printf("%s sent from %s to %s\n", amount, fromAccountInput, toAccountInput)
	return nil
}
