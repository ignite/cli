package starportcmd

import "github.com/spf13/cobra"

func NewNetworkProposal() *cobra.Command {
	c := &cobra.Command{
		Use:               "proposal",
		Short:             "Proposals related to starting network",
		PersistentPreRunE: ensureSPNAccountHook,
	}
	c.AddCommand(NewNetworkProposalList())
	c.AddCommand(NewNetworkProposalDescribe())
	c.AddCommand(NewNetworkProposalApprove())
	c.AddCommand(NewNetworkProposalReject())
	c.AddCommand(NewNetworkProposalTest())
	return c
}
