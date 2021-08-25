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
	flagAdvanced            = "advanced"
	flagSourceRPC           = "source-rpc"
	flagTargetRPC           = "target-rpc"
	flagSourceFaucet        = "source-faucet"
	flagTargetFaucet        = "target-faucet"
	flagSourcePort          = "source-port"
	flagSourceVersion       = "source-version"
	flagTargetPort          = "target-port"
	flagTargetVersion       = "target-version"
	flagSourceGasPrice      = "source-gasprice"
	flagTargetGasPrice      = "target-gasprice"
	flagSourceGasLimit      = "source-gaslimit"
	flagTargetGasLimit      = "target-gaslimit"
	flagSourceAddressPrefix = "source-prefix"
	flagTargetAddressPrefix = "target-prefix"
	flagOrdered             = "ordered"

	relayerSource = "source"
	relayerTarget = "target"

	defaultSourceRPCAddress = "http://localhost:26657"
	defaultTargetRPCAddress = "https://rpc.cosmos.network:443"

	defautSourceGasPrice      = "0.00025stake"
	defautTargetGasPrice      = "0.025uatom"
	defautSourceGasLimit      = 300000
	defautTargetGasLimit      = 300000
	defautSourceAddressPrefix = "cosmos"
	defautTargetAddressPrefix = "cosmos"
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
	c.Flags().BoolP(flagAdvanced, "a", false, "Advanced configuration options for custom IBC modules")
	c.Flags().String(flagSourceRPC, "", "RPC address of the source chain")
	c.Flags().String(flagTargetRPC, "", "RPC address of the target chain")
	c.Flags().String(flagSourceFaucet, "", "Faucet address of the source chain")
	c.Flags().String(flagTargetFaucet, "", "Faucet address of the target chain")
	c.Flags().String(flagSourcePort, "", "IBC port ID on the source chain")
	c.Flags().String(flagSourceVersion, "", "Module version on the source chain")
	c.Flags().String(flagTargetPort, "", "IBC port ID on the target chain")
	c.Flags().String(flagTargetVersion, "", "Module version on the target chain")
	c.Flags().String(flagSourceGasPrice, "", "Gas price used for transactions on source chain")
	c.Flags().String(flagTargetGasPrice, "", "Gas price used for transactions on target chain")
	c.Flags().Int64(flagSourceGasLimit, 0, "Gas limit used for transactions on source chain")
	c.Flags().Int64(flagTargetGasLimit, 0, "Gas limit used for transactions on target chain")
	c.Flags().String(flagSourceAddressPrefix, "", "Address prefix of the source chain")
	c.Flags().String(flagTargetAddressPrefix, "", "Address prefix of the target chain")
	c.Flags().Bool(flagOrdered, false, "Set the channel as ordered")

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
		sourceGasLimit      int64
		targetGasLimit      int64
		sourceAddressPrefix string
		targetAddressPrefix string
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
		questionSourceGasLimit = cliquiz.NewQuestion(
			"Source Gas Limit",
			&sourceGasLimit,
			cliquiz.DefaultAnswer(defautSourceGasLimit),
			cliquiz.Required(),
		)
		questionTargetGasLimit = cliquiz.NewQuestion(
			"Target Gas Limit",
			&targetGasLimit,
			cliquiz.DefaultAnswer(defautTargetGasLimit),
			cliquiz.Required(),
		)
		questionSourceAddressPrefix = cliquiz.NewQuestion(
			"Source Address Prefix",
			&sourceAddressPrefix,
			cliquiz.DefaultAnswer(defautSourceAddressPrefix),
			cliquiz.Required(),
		)
		questionTargetAddressPrefix = cliquiz.NewQuestion(
			"Target Address Prefix",
			&targetAddressPrefix,
			cliquiz.DefaultAnswer(defautTargetAddressPrefix),
			cliquiz.Required(),
		)
	)

	// Get flags
	advanced, err := cmd.Flags().GetBool(flagAdvanced)
	if err != nil {
		return err
	}
	sourceRPCAddress, err = cmd.Flags().GetString(flagSourceRPC)
	if err != nil {
		return err
	}
	sourceFaucetAddress, err = cmd.Flags().GetString(flagSourceFaucet)
	if err != nil {
		return err
	}
	targetRPCAddress, err = cmd.Flags().GetString(flagTargetRPC)
	if err != nil {
		return err
	}
	targetFaucetAddress, err = cmd.Flags().GetString(flagTargetFaucet)
	if err != nil {
		return err
	}
	sourcePort, err = cmd.Flags().GetString(flagSourcePort)
	if err != nil {
		return err
	}
	sourceVersion, err = cmd.Flags().GetString(flagSourceVersion)
	if err != nil {
		return err
	}
	targetPort, err = cmd.Flags().GetString(flagTargetPort)
	if err != nil {
		return err
	}
	targetVersion, err = cmd.Flags().GetString(flagTargetVersion)
	if err != nil {
		return err
	}
	sourceGasPrice, err = cmd.Flags().GetString(flagSourceGasPrice)
	if err != nil {
		return err
	}
	targetGasPrice, err = cmd.Flags().GetString(flagTargetGasPrice)
	if err != nil {
		return err
	}
	sourceGasLimit, err = cmd.Flags().GetInt64(flagSourceGasLimit)
	if err != nil {
		return err
	}
	targetGasLimit, err = cmd.Flags().GetInt64(flagTargetGasLimit)
	if err != nil {
		return err
	}
	sourceAddressPrefix, err = cmd.Flags().GetString(flagSourceAddressPrefix)
	if err != nil {
		return err
	}
	targetAddressPrefix, err = cmd.Flags().GetString(flagTargetAddressPrefix)
	if err != nil {
		return err
	}
	ordered, err := cmd.Flags().GetBool(flagOrdered)
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
	if sourceGasLimit == 0 {
		questions = append(questions, questionSourceGasLimit)
	}
	if targetGasLimit == 0 {
		questions = append(questions, questionTargetGasLimit)
	}
	if sourceAddressPrefix == "" {
		questions = append(questions, questionSourceAddressPrefix)
	}
	if targetAddressPrefix == "" {
		questions = append(questions, questionTargetAddressPrefix)
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
		sourceGasLimit,
		sourceAddressPrefix,
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
		targetGasLimit,
		targetAddressPrefix,
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
	path, err := sourceChain.Connect(cmd.Context(), targetChain, channelOptions...)
	if err != nil {
		return err
	}

	s.Stop()

	info, err := xrelayer.Info(cmd.Context())
	if err != nil {
		return err
	}

	fmt.Printf("⛓  Configured chains: %s\n\n", color.Green.Sprint(path.ID))
	fmt.Printf(`Note: mnemonics for relayer accounts are stored in %s unencrypted.
This may change in the future. Until then, use them only for small amounts of tokens.
`, info.ConfigPath)

	return nil
}

// initChain initializes chain information for the relayer connection
func initChain(
	cmd *cobra.Command,
	s *clispinner.Spinner,
	name,
	rpcAddr,
	faucetAddr,
	gasPrice string,
	gasLimit int64,
	addressPrefix string,
) (*xrelayer.Chain, error) {
	defer s.Stop()
	s.SetText("Initializing chain...").Start()

	c, err := xrelayer.NewChain(
		cmd.Context(),
		rpcAddr,
		xrelayer.WithFaucet(faucetAddr),
		xrelayer.WithGasPrice(gasPrice),
		xrelayer.WithGasLimit(gasLimit),
		xrelayer.WithAddressPrefix(addressPrefix),
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
