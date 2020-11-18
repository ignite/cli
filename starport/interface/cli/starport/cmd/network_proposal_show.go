package starportcmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

func NewNetworkProposalDescribe() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [chain-id] [number]",
		Short: "Show details of a proposal",
		RunE:  networkProposalDescribeHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkProposalDescribeHandler(cmd *cobra.Command, args []string) error {
	b, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return err
	}

	proposal, err := b.ProposalGet(context.Background(), args[0], int(id))
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
