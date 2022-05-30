package ignitecmd

import (
	"github.com/spf13/cobra"
)

var rpcAddress string

const (
	flagRPC         = "rpc"
	rpcAddressLocal = "tcp://localhost:26657"
)

func NewNode() *cobra.Command {
	c := &cobra.Command{
		Use:   "node [command]",
		Short: "Make calls to a live blockchain node",
		Args:  cobra.ExactArgs(1),
	}

	c.PersistentFlags().StringVar(&rpcAddress, flagRPC, rpcAddressLocal, "<host>:<port> to tendermint rpc interface for this chain")

	c.AddCommand(NewNodeQuery())
	c.AddCommand(NewNodeTx())

	return c
}

func getRPC(cmd *cobra.Command) (rpc string) {
	rpc, _ = cmd.Flags().GetString(flagRPC)
	return
}
