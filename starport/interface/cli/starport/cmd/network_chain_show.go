package starportcmd

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

func NewNetworkChainShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [chain-id]",
		Short: "Show details of a chain",
		RunE:  networkChainShowHandler,
		Args:  cobra.ExactArgs(1),
	}
	return c
}

func networkChainShowHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	chain, err := nb.ChainShow(context.Background(), args[0])
	if err != nil {
		return err
	}
	chainyaml, err := yaml.Marshal(chain)
	if err != nil {
		return err
	}
	fmt.Println(string(chainyaml))

	return nil
}
