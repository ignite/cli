package starportcmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/numbers"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

const (
	flagDebug = "debug"
)

func NewNetworkProposalVerify() *cobra.Command {
	c := &cobra.Command{
		Use:   "verify [chain-id] [number<,...>]",
		Short: "Simulate and verify proposals validity",
		RunE:  networkProposalVerifyHandler,
		Args:  cobra.ExactArgs(2),
	}
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().Bool(flagDebug, false, "show output of verification command in the console")
	return c
}

func networkProposalVerifyHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New()
	defer s.Stop()

	var (
		chainID      = args[0]
		proposalList = args[1]
	)

	ev := events.NewBus()
	go printEvents(ev, s)

	// Initialize the blockchain
	nb, err := newNetworkBuilder(cmd.Context(), networkbuilder.CollectEvents(ev))
	if err != nil {
		return err
	}

	ids, err := numbers.ParseList(proposalList)
	if err != nil {
		return err
	}

	s.SetText("Verifying proposals...")
	s.Start()

	// Check verbose flag
	out := ioutil.Discard
	debugSet, err := cmd.Flags().GetBool(flagDebug)
	if err != nil {
		return err
	}
	if debugSet {
		out = os.Stdout
	}

	err = nb.VerifyProposals(cmd.Context(), chainID, ids, out)
	s.Stop()
	var verificationError networkbuilder.VerificationError
	if errors.As(err, &verificationError) {
		fmt.Printf("Proposal(s) %s invalid 🔍❌️\nError: %s️\n",
			numbers.List(ids, "#"),
			err.Error(),
		)
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("Proposal(s) %s verified 🔍✅️\n", numbers.List(ids, "#"))

	return nil
}
