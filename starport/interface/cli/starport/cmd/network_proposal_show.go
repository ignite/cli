package starportcmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// NewNetworkProposalShow creates a new command to show a proposal in a chain.
func NewNetworkProposalShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [chain-id] [number]",
		Short: "Show details of a proposal",
		RunE:  networkProposalShowHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkProposalShowHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return err
	}

	proposal, err := nb.ProposalGet(context.Background(), args[0], int(id))
	if err != nil {
		return err
	}
	proposalyaml, err := yaml.Marshal(proposal)
	if err != nil {
		return err
	}
	fmt.Println(string(proposalyaml))

	return nil
}
