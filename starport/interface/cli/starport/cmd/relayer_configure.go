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
	advancedFlag       = "advanced"
	sourceRPCFlag      = "source-rpc"
	targetRPCFlag      = "target-rpc"
	sourceFaucetFlag   = "source-faucet"
	targetFaucetFlag   = "target-faucet"
	sourcePortFlag     = "source-port"
	sourceVersionFlag  = "source-version"
	targetPortFlag     = "target-port"
	targetVersionFlag  = "target-version"
	sourceGasPriceFlag = "source-gasprice"
	targetGasPriceFlag = "target-gasprice"
	orderedFlag        = "ordered"

	relayerSource = "source"
	relayerTarget = "target"

	defaultSourceRPCAddress = "http://localhost:26657"
	defaultTargetRPCAddress = "https://rpc.cosmos.network:443"

	defautSourceGasPrice = "0.025stake"
	defautTargetGasPrice = "0.025uatom"
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
	c.Flags().String(sourceRPCFlag, "", "RPC address of the source chain")
	c.Flags().String(targetRPCFlag, "", "RPC address of the target chain")
	c.Flags().String(sourceFaucetFlag, "", "Faucet address of the source chain")
	c.Flags().String(targetFaucetFlag, "", "Faucet address of the target chain")
	c.Flags().String(sourcePortFlag, "", "IBC port ID on the source chain")
	c.Flags().String(sourceVersionFlag, "", "Module version on the source chain")
	c.Flags().String(targetPortFlag, "", "IBC port ID on the target chain")
	c.Flags().String(targetVersionFlag, "", "Module version on the target chain")
	c.Flags().String(sourceGasPriceFlag, "", "Gas price used for transactions on source chain")
	c.Flags().String(targetGasPriceFlag, "", "Gas price used for transactions on target chain")
	c.Flags().Bool(orderedFlag, false, "Set the channel as ordered")

	return c
}

func relayerConfigureHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().Stop()
	defer s.Stop()

	printSection("Setting up chains")

	// basic configuration
	var (
		sourceRPCAddress    string
		targetRPCAddress    string
		sourceFaucetAddress string
		targetFaucetAddress string
		sourceGasPrice      string
		targetGasPrice      string
	)

	// advanced configuration for the channel
	var (
		sourcePort    string
		sourceVersion string
		targetPort    string
		targetVersion string
	)

	// questions
	var (
		questionSourceRPCAddress = cliquiz.NewQuestion(
			"Source RPC",
			&sourceRPCAddress,
			cliquiz.DefaultAnswer(defaultSourceRPCAddress),
			cliquiz.Required(),
		)
		questionSourceFaucet = cliquiz.NewQuestion(
			"Source Faucet",
			&sourceFaucetAddress,
		)
		questionTargetRPCAddress = cliquiz.NewQuestion(
			"Target RPC",
			&targetRPCAddress,
			cliquiz.DefaultAnswer(defaultTargetRPCAddress),
			cliquiz.Required(),
		)
		questionTargetFaucet = cliquiz.NewQuestion(
			"Target Faucet",
			&targetFaucetAddress,
		)
		questionSourcePort = cliquiz.NewQuestion(
			"Source Port",
			&sourcePort,
			cliquiz.DefaultAnswer(xrelayer.TransferPort),
			cliquiz.Required(),
		)
		questionSourceVersion = cliquiz.NewQuestion(
			"Source Version",
			&sourceVersion,
			cliquiz.DefaultAnswer(xrelayer.TransferVersion),
			cliquiz.Required(),
		)
		questionTargetPort = cliquiz.NewQuestion(
			"Target Port",
			&targetPort,
			cliquiz.DefaultAnswer(xrelayer.TransferPort),
			cliquiz.Required(),
		)
		questionTargetVersion = cliquiz.NewQuestion(
			"Target Version",
			&targetVersion,
			cliquiz.DefaultAnswer(xrelayer.TransferVersion),
			cliquiz.Required(),
		)
		questionSourceGasPrice = cliquiz.NewQuestion(
			"Source Gas Price",
			&sourceGasPrice,
			cliquiz.DefaultAnswer(defautSourceGasPrice),
			cliquiz.Required(),
		)
		questionTargetGasPrice = cliquiz.NewQuestion(
			"Target Gas Price",
			&targetGasPrice,
			cliquiz.DefaultAnswer(defautTargetGasPrice),
			cliquiz.Required(),
		)
	)

	// Get flags
	advanced, err := cmd.Flags().GetBool(advancedFlag)
	if err != nil {
		return err
	}
	sourceRPCAddress, err = cmd.Flags().GetString(sourceRPCFlag)
	if err != nil {
		return err
	}
	sourceFaucetAddress, err = cmd.Flags().GetString(sourceFaucetFlag)
	if err != nil {
		return err
	}
	targetRPCAddress, err = cmd.Flags().GetString(targetRPCFlag)
	if err != nil {
		return err
	}
	targetFaucetAddress, err = cmd.Flags().GetString(targetFaucetFlag)
	if err != nil {
		return err
	}
	sourcePort, err = cmd.Flags().GetString(sourcePortFlag)
	if err != nil {
		return err
	}
	sourceVersion, err = cmd.Flags().GetString(sourceVersionFlag)
	if err != nil {
		return err
	}
	targetPort, err = cmd.Flags().GetString(targetPortFlag)
	if err != nil {
		return err
	}
	targetVersion, err = cmd.Flags().GetString(targetVersionFlag)
	if err != nil {
		return err
	}
	sourceGasPrice, err = cmd.Flags().GetString(sourceGasPriceFlag)
	if err != nil {
		return err
	}
	targetGasPrice, err = cmd.Flags().GetString(targetGasPriceFlag)
	if err != nil {
		return err
	}
	ordered, err := cmd.Flags().GetBool(orderedFlag)
	if err != nil {
		return err
	}

	var questions []cliquiz.Question

	// get information from prompt if flag not provided
	if sourceRPCAddress == "" {
		questions = append(questions, questionSourceRPCAddress)
	}
	if sourceFaucetAddress == "" {
		questions = append(questions, questionSourceFaucet)
	}
	if targetRPCAddress == "" {
		questions = append(questions, questionTargetRPCAddress)
	}
	if targetFaucetAddress == "" {
		questions = append(questions, questionTargetFaucet)
	}
	if sourceGasPrice == "" {
		questions = append(questions, questionSourceGasPrice)
	}
	if targetGasPrice == "" {
		questions = append(questions, questionTargetGasPrice)
	}

	// advanced information
	if advanced {
		if sourcePort == "" {
			questions = append(questions, questionSourcePort)
		}
		if sourceVersion == "" {
			questions = append(questions, questionSourceVersion)
		}
		if targetPort == "" {
			questions = append(questions, questionTargetPort)
		}
		if targetVersion == "" {
			questions = append(questions, questionTargetVersion)
		}
	}

	if len(questions) > 0 {
		if err := cliquiz.Ask(questions...); err != nil {
			return err
		}
	}

	fmt.Println()
	s.SetText("Fetching chain info...")

	// initialize the chains
	sourceChain, err := initChain(
		cmd,
		s,
		relayerSource,
		sourceRPCAddress,
		sourceFaucetAddress,
		sourceGasPrice,
	)
	if err != nil {
		return err
	}

	targetChain, err := initChain(
		cmd,
		s,
		relayerTarget,
		targetRPCAddress,
		targetFaucetAddress,
		targetGasPrice,
	)
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

		if ordered {
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

// initChain initializes chain information for the relayer connection
func initChain(cmd *cobra.Command, s *clispinner.Spinner, name, rpcAddr, faucetAddr, gasPrice string) (*xrelayer.Chain, error) {
	defer s.Stop()
	s.SetText("Initializing chain...").Start()

	c, err := xrelayer.NewChain(
		cmd.Context(),
		rpcAddr,
		xrelayer.WithFaucet(faucetAddr),
		xrelayer.WithGasPrice(gasPrice),
	)
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
	balance := coins.String()
	if balance == "" {
		balance = "-"
	}
	fmt.Printf(" |· (balance: %s)\n\n", balance)

	return c, nil
}
