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
	advancedFlag = "advanced"

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
	c.Flags().BoolP(advancedFlag, "a", false, "Advanced configuration options for custom IBC modules")
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

	// advanced configuration for the channel
	var (
		sourcePort          string
		sourceVersion    string
		targetPort string
		targetVersion string
		ordered string
	)

	// check if advanced configuration
	advanced, err := cmd.Flags().GetBool(advancedFlag)
	if err != nil {
		return err
	}

	var questions []cliquiz.Question

	// source configuration
	questions = append(questions,
		cliquiz.NewQuestion("Source RPC",
			&sourceRPCAddress,
			cliquiz.DefaultAnswer(defaultSourceRPCAddress),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Source Faucet",
			&sourceFaucetAddress,
		),
	)

	// advanced source configuration
	if advanced {
		questions = append(questions,
			cliquiz.NewQuestion("Source Port",
				&sourcePort,
				cliquiz.DefaultAnswer(xrelayer.TransferPort),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Source Version",
				&sourceVersion,
				cliquiz.DefaultAnswer(xrelayer.TransferVersion),
				cliquiz.Required(),
			),
			)
	}

	// target configuration
	questions = append(questions,
		cliquiz.NewQuestion("Target RPC",
			&targetRPCAddress,
			cliquiz.DefaultAnswer(defaultTargetRPCAddress),
			cliquiz.Required(),
		),
		cliquiz.NewQuestion("Target Faucet",
			&targetFaucetAddress,
		),
	)

	// advanced target configuration and other advanced configuration
	if advanced {
		questions = append(questions,
			cliquiz.NewQuestion("Target Port",
				&targetPort,
				cliquiz.DefaultAnswer(xrelayer.TransferPort),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Target Version",
				&targetVersion,
				cliquiz.DefaultAnswer(xrelayer.TransferVersion),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Ordered Channel (yes/no)",
				&ordered,
				cliquiz.DefaultAnswer("no"),
				cliquiz.Required(),
			),
		)
	}

	if err := cliquiz.Ask(questions...); err != nil {
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

	// initialize the chains
	sourceChain, err := init(relayerSource, sourceRPCAddress, sourceFaucetAddress)
	if err != nil {
		return err
	}

	targetChain, err := init(relayerTarget, targetRPCAddress, targetFaucetAddress)
	if err != nil {
		return err
	}

	s.SetText("Configuring...").Start()

	// sets advanced channel options
	var channelOptions []xrelayer.ChannelOption
	if advanced {
		channelOptions = append(channelOptions,
			xrelayer.SourcePort(sourcePort),
			xrelayer.SourceVersion(sourceVersion),
			xrelayer.TargetPort(targetPort),
			xrelayer.TargetVersion(targetVersion),
		)

		if ordered == "yes" {
			channelOptions = append(channelOptions, xrelayer.Ordered())
		}
	}

	// create the connection configuration
	connectionID, err := sourceChain.Connect(cmd.Context(), targetChain, channelOptions...)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("⛓  Configured chains: %s\n\n", color.Green.Sprint(connectionID))

	return nil
}
