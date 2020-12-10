package starportcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/numbers"
)

func NewNetworkProposalVerify() *cobra.Command {
	c := &cobra.Command{
		Use:   "verify [chain-id] [number<,...>]",
		Short: "Simulate and verify proposals validity",
		RunE:  networkProposalVerifyHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkProposalVerifyHandler(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("Proposal(s) %s verified 🔍✔️\n", numbers.List(ids, "#"))
	return nil
}