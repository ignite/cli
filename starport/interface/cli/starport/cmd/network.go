package starportcmd

import (
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

var spnAddress string

func NewNetwork() *cobra.Command {
	c := &cobra.Command{
		Use:   "network",
		Short: "Create and start Blochains collaboratively",
		Args:  cobra.ExactArgs(1),
	}

	// configure flags.
	c.PersistentFlags().StringVarP(&spnAddress, "spn-address", "s", "localhost:26657", "An SPN node address")

	// add sub commands.
	c.AddCommand(NewNetworkChain())
	c.AddCommand(NewNetworkAccount())
	return c
}

func newNetworkBuilder(options ...networkbuilder.Option) (*networkbuilder.Builder, error) {
	var spnoptions []spn.Option
	// use test keyring backend on Gitpod in order to prevent prompting for keyring
	// password. This happens because Gitpod uses containers.
	//
	// when not on Gitpod, OS keyring backend is used which only asks password once.
	if os.Getenv("GITPOD_WORKSPACE_ID") != "" {
		spnoptions = append(spnoptions, spn.Keyring(keyring.BackendTest))
	}
	spnclient, err := spn.New(spnAddress, spnoptions...)
	if err != nil {
		return nil, err
	}
	return networkbuilder.New(spnclient, options...)
}
