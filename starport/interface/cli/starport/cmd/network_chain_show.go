package starportcmd

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// NewNetworkChainShow creates a new chain show command to show
// a chain on SPN.
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

	chainID := args[0]

	c, err := nb.ShowChain(context.Background(), chainID)
	if err != nil {
		return err
	}
	chainyaml, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	info, err := nb.LaunchInformation(context.Background(), chainID)
	if err != nil {
		return err
	}
	infoyaml, err := yaml.Marshal(info)
	if err != nil {
		return err
	}

	fmt.Printf("\nChain:\n---\n%s\n\nLaunch Information:\n---\n%s", string(chainyaml), string(infoyaml))

	return nil
}
