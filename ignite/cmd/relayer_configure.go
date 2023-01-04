package ignitecmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite/cli/ignite/pkg/cliui/entrywriter"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/relayer"
	relayerconfig "github.com/ignite/cli/ignite/pkg/relayer/config"
)

const (
	flagAdvanced            = "advanced"
	flagSourceAccount       = "source-account"
	flagTargetAccount       = "target-account"
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
	flagReset               = "reset"
	flagSourceClientID      = "source-client-id"
	flagTargetClientID      = "target-client-id"

	RelayerSource = "source"
	RelayerTarget = "target"

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

	c.Flags().BoolP(flagAdvanced, "a", false, "advanced configuration options for custom IBC modules")
	c.Flags().String(flagSourceRPC, "", "RPC address of the source chain")
	c.Flags().String(flagTargetRPC, "", "RPC address of the target chain")
	c.Flags().String(flagSourceFaucet, "", "faucet address of the source chain")
	c.Flags().String(flagTargetFaucet, "", "faucet address of the target chain")
	c.Flags().String(flagSourcePort, "", "IBC port ID on the source chain")
	c.Flags().String(flagSourceVersion, "", "module version on the source chain")
	c.Flags().String(flagTargetPort, "", "IBC port ID on the target chain")
	c.Flags().String(flagTargetVersion, "", "module version on the target chain")
	c.Flags().String(flagSourceGasPrice, "", "gas price used for transactions on source chain")
	c.Flags().String(flagTargetGasPrice, "", "gas price used for transactions on target chain")
	c.Flags().Int64(flagSourceGasLimit, 0, "gas limit used for transactions on source chain")
	c.Flags().Int64(flagTargetGasLimit, 0, "gas limit used for transactions on target chain")
	c.Flags().String(flagSourceAddressPrefix, "", "address prefix of the source chain")
	c.Flags().String(flagTargetAddressPrefix, "", "address prefix of the target chain")
	c.Flags().String(flagSourceAccount, "", "source Account")
	c.Flags().String(flagTargetAccount, "", "target Account")
	c.Flags().Bool(flagOrdered, false, "set the channel as ordered")
	c.Flags().BoolP(flagReset, "r", false, "reset the relayer config")
	c.Flags().String(flagSourceClientID, "", "use a custom client id for source")
	c.Flags().String(flagTargetClientID, "", "use a custom client id for target")
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())

	return c
}

func relayerConfigureHandler(cmd *cobra.Command, _ []string) (err error) {
	defer func() {
		err = handleRelayerAccountErr(err)
	}()

	session := cliui.New()
	defer session.End()

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
		cosmosaccount.WithHome(getKeyringDir(cmd)),
	)
	if err != nil {
		return err
	}

	if err := ca.EnsureDefaultAccount(); err != nil {
		return err
	}

	if err := printSection(session, "Setting up chains"); err != nil {
		return err
	}

	// basic configuration
	var (
		sourceAccount       string
		targetAccount       string
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
		questionSourceAccount = cliquiz.NewQuestion(
			"Source Account",
			&sourceAccount,
			cliquiz.DefaultAnswer(cosmosaccount.DefaultAccount),
			cliquiz.Required(),
		)
		questionTargetAccount = cliquiz.NewQuestion(
			"Target Account",
			&targetAccount,
			cliquiz.DefaultAnswer(cosmosaccount.DefaultAccount),
			cliquiz.Required(),
		)
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
			cliquiz.DefaultAnswer(relayer.TransferPort),
			cliquiz.Required(),
		)
		questionSourceVersion = cliquiz.NewQuestion(
			"Source Version",
			&sourceVersion,
			cliquiz.DefaultAnswer(relayer.TransferVersion),
			cliquiz.Required(),
		)
		questionTargetPort = cliquiz.NewQuestion(
			"Target Port",
			&targetPort,
			cliquiz.DefaultAnswer(relayer.TransferPort),
			cliquiz.Required(),
		)
		questionTargetVersion = cliquiz.NewQuestion(
			"Target Version",
			&targetVersion,
			cliquiz.DefaultAnswer(relayer.TransferVersion),
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
	sourceAccount, err = cmd.Flags().GetString(flagSourceAccount)
	if err != nil {
		return err
	}
	targetAccount, err = cmd.Flags().GetString(flagTargetAccount)
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
	var (
		sourceClientID, _ = cmd.Flags().GetString(flagSourceClientID)
		targetClientID, _ = cmd.Flags().GetString(flagTargetClientID)
		reset, _          = cmd.Flags().GetBool(flagReset)

		questions []cliquiz.Question
	)

	// get information from prompt if flag not provided
	if sourceAccount == "" {
		questions = append(questions, questionSourceAccount)
	}
	if targetAccount == "" {
		questions = append(questions, questionTargetAccount)
	}
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

	session.PauseSpinner()
	if len(questions) > 0 {
		if err := session.Ask(questions...); err != nil {
			return err
		}
	}

	if reset {
		if err := relayerconfig.Delete(); err != nil {
			return err
		}
	}

	session.StartSpinner("Fetching chain info...")

	session.Println()
	r := relayer.New(ca)

	// initialize the chains
	sourceChain, err := InitChain(
		cmd,
		r,
		session,
		RelayerSource,
		sourceAccount,
		sourceRPCAddress,
		sourceFaucetAddress,
		sourceGasPrice,
		sourceGasLimit,
		sourceAddressPrefix,
		sourceClientID,
	)
	if err != nil {
		return err
	}

	if err := sourceChain.EnsureChainSetup(cmd.Context()); err != nil {
		return err
	}

	targetChain, err := InitChain(
		cmd,
		r,
		session,
		RelayerTarget,
		targetAccount,
		targetRPCAddress,
		targetFaucetAddress,
		targetGasPrice,
		targetGasLimit,
		targetAddressPrefix,
		targetClientID,
	)
	if err != nil {
		return err
	}

	if err := targetChain.EnsureChainSetup(cmd.Context()); err != nil {
		return err
	}

	session.StartSpinner("Configuring...")

	// sets advanced channel options
	var channelOptions []relayer.ChannelOption
	if advanced {
		channelOptions = append(channelOptions,
			relayer.SourcePort(sourcePort),
			relayer.SourceVersion(sourceVersion),
			relayer.TargetPort(targetPort),
			relayer.TargetVersion(targetVersion),
		)

		if ordered {
			channelOptions = append(channelOptions, relayer.Ordered())
		}
	}

	// create the connection configuration
	id, err := sourceChain.Connect(targetChain, channelOptions...)
	if err != nil {
		return err
	}

	return session.Printf("⛓  Configured chains: %s\n\n", color.Green.Sprint(id))
}

// InitChain initializes chain information for the relayer connection.
func InitChain(
	cmd *cobra.Command,
	r relayer.Relayer,
	session *cliui.Session,
	name,
	accountName,
	rpcAddr,
	faucetAddr,
	gasPrice string,
	gasLimit int64,
	addressPrefix,
	clientID string,
) (*relayer.Chain, error) {
	defer session.StopSpinner()
	session.StartSpinner(fmt.Sprintf("Initializing chain %s...", name))

	c, account, err := r.NewChain(
		accountName,
		rpcAddr,
		relayer.WithFaucet(faucetAddr),
		relayer.WithGasPrice(gasPrice),
		relayer.WithGasLimit(gasLimit),
		relayer.WithAddressPrefix(addressPrefix),
		relayer.WithClientID(clientID),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot resolve %s", name)
	}

	accountAddr, err := account.Address(addressPrefix)
	if err != nil {
		return nil, err
	}

	session.StopSpinner()
	session.Printf("🔐  Account on %q is %s(%s)\n \n", name, accountName, accountAddr)
	session.StartSpinner(color.Yellow.Sprintf("trying to receive tokens from a faucet..."))

	coins, err := c.TryRetrieve(cmd.Context())

	session.StopSpinner()
	session.Print(" |· ")
	if err != nil {
		session.Println(color.Yellow.Sprintf(err.Error()))
	} else {
		session.Println(color.Green.Sprintf("received coins from a faucet"))
	}

	balance := coins.String()
	if balance == "" {
		balance = entrywriter.None
	}
	session.Printf(" |· (balance: %s)\n\n", balance)

	return c, nil
}
