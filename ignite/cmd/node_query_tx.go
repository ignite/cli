package ignitecmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
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
	bz, err := hex.DecodeString(args[0])
	if err != nil {
		return err
	}
	rpc, err := sdkclient.NewClientFromNode(getRPC(cmd))
	if err != nil {
		return err
	}
	resp, err := rpc.Tx(cmd.Context(), bz, false)
	if err != nil {
		return err
	}
	bz, err = json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bz))
	return nil
}
