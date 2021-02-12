package starportcmd

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/xrelayer"
)

const (
	relayerSource = "source"
	relayerTarget = "target"

	defaultSourceRPCAddress = "http://localhost:26657"
	defaultTargetRPCAddress = "https://rpc.alpha.starport.network:443"
)

// NewRelayerConfigure returns a new relayer configure command.
// faucet addresses are optional and connect command will try to guess the address
// when not provided. even if auto retrieving coins fails, connect command will complete with success.
func NewRelayerConfigure() *cobra.Command {
	c := &cobra.Command{
		Use:     "configure",
		Short:   "Configure source and target chains for relaying",
		Aliases: []string{"conf"},
		RunE:    relayerConfigureHandler,
	}
	return c
}

func relayerConfigureHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().Stop()
	defer s.Stop()

	printSection("Setting up chains")

	var (
		sourceRPCAddress    string
		targetRPCAddress    string
		sourceFaucetAddress string
		targetFaucetAddress string
	)

	if err := cliquiz.Ask(
		cliquiz.NewQuestion("Source RPC",
			&sourceRPCAddress,
			cliquiz.DefaultAnswer(defaultSourceRPCAddress),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Source Faucet",
			&sourceFaucetAddress,
		),
		cliquiz.NewQuestion("Target RPC",
			&targetRPCAddress,
			cliquiz.DefaultAnswer(defaultTargetRPCAddress),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Target Faucet",
			&targetFaucetAddress,
		),
	); err != nil {
		return err
	}

	fmt.Println()
	s.SetText("Fetching chain info...")

	init := func(name, rpcAddr, faucetAddr string) (*xrelayer.Chain, error) {
		defer s.Stop()
		s.SetText("Initializing chain...").Start()

		c, err := xrelayer.NewChain(cmd.Context(), rpcAddr, xrelayer.WithFaucet(faucetAddr))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot resolve %s", name)
		}

		account, err := c.Account(cmd.Context())
		if err != nil {
			return nil, err
		}

		s.Stop()

		fmt.Printf("🔐  Account on %q is %q\n \n", name, account.Address)
		s.
			SetCharset(spinner.CharSets[9]).
			SetColor("white").
			SetPrefix(" |·").
			SetText(color.Yellow.Sprintf("trying to receive tokens from a faucet...")).
			Start()

		err = c.TryFaucet(cmd.Context())
		s.Stop()

		fmt.Print(" |· ")
		if err != nil {
			fmt.Println(color.Yellow.Sprintf(err.Error()))
		} else {
			fmt.Println(color.Green.Sprintf("received coins from a faucet"))
		}

		coins, err := c.Balance(cmd.Context())
		if err != nil {
			return nil, err
		}
		fmt.Printf(" |· (balance: %s)\n\n", coins)

		return c, nil
	}

	sourceChain, err := init(relayerSource, sourceRPCAddress, sourceFaucetAddress)
	if err != nil {
		return err
	}

	targetChain, err := init(relayerTarget, targetRPCAddress, targetFaucetAddress)
	if err != nil {
		return err
	}

	s.SetText("Configuring...").Start()

	connectionID, err := sourceChain.Connect(cmd.Context(), targetChain)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("⛓  Configured chains: %s\n\n", color.Green.Sprint(connectionID))

	return nil
}
