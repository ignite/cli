package ignitecmd

import (
	"encoding/hex"
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
)

func NewNodeQueryTx() *cobra.Command {
	c := &cobra.Command{
		Use:   "tx [hash]",
		Short: "Query for transaction by hash",
		RunE:  nodeQueryTxHandler,
		Args:  cobra.ExactArgs(1),
	}
	return c
}

func nodeQueryTxHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusQuerying))
	defer session.End()

	bz, err := hex.DecodeString(args[0])
	if err != nil {
		return err
	}
	client, err := newNodeCosmosClient(cmd)
	if err != nil {
		return err
	}

	resp, err := client.RPC.Tx(cmd.Context(), bz, false)
	if err != nil {
		return err
	}
	bz, err = json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}

	return session.Println(string(bz))
}
