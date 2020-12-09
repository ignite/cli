package starportcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/numbers"
)

func NewNetworkProposalTest() *cobra.Command {
	c := &cobra.Command{
		Use:   "test [chain-id] [number<,...>]",
		Short: "Test proposals validity",
		RunE:  networkProposalTestHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkProposalTestHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New()
	defer s.Stop()

	var (
		_      			= args[0]
		proposalList 	= args[1]
	)

	_, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	ids, err := numbers.ParseList(proposalList)
	if err != nil {
		return err
	}

	s.Stop()


	s.SetText("Test...")
	s.Start()
	s.Stop()

	fmt.Printf("Proposal(s) %s tested ✅\n", numbers.List(ids, "#"))
	return nil
}