package starportcmd

import "github.com/spf13/cobra"

// NewNetworkProposal creates a proposal command that holds some other
// sub commands.
func NewNetworkProposal() *cobra.Command {
	c := &cobra.Command{
		Use:               "proposal",
		Short:             "Proposals related to starting network",
		PersistentPreRunE: ensureSPNAccountHook,
	}
	c.AddCommand(NewNetworkProposalList())
	c.AddCommand(NewNetworkProposalShow())
	c.AddCommand(NewNetworkProposalApprove())
	c.AddCommand(NewNetworkProposalReject())
	c.AddCommand(NewNetworkProposalVerify())
	return c
}
