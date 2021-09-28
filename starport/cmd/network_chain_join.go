package starportcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/networkbuilder"
	"github.com/tendermint/tendermint/libs/os"
)

const (
	flagGentx = "gentx"
	flagPeer  = "peer"
)

const (
	flagKeyringBackend = "keyring-backend"
)

// NewNetworkChainJoin creates a new chain join command to join
// to a network as a validator.
func NewNetworkChainJoin() *cobra.Command {
	c := &cobra.Command{
		Use:   "join [chain-id]",
		Short: "Propose to join to a network as a validator",
		RunE:  networkChainJoinHandler,
		Args:  cobra.ExactArgs(1),
	}
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().String(flagGentx, "", "Path to a gentx file (optional)")
	c.Flags().String(flagPeer, "", "Configure peer in node-id@host:port format (optional)")
	c.Flags().String(flagKeyringBackend, "os", "Keyring backend used for the blockchain account")
	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	chainID := args[0]
	gentxPath, _ := cmd.Flags().GetString(flagGentx)
	publicAddress, _ := cmd.Flags().GetString(flagPeer)

	s := clispinner.New()
	defer s.Stop()

	ev := events.NewBus()
	go printEvents(ev, s)

	nb, err := newNetworkBuilder(cmd.Context(), networkbuilder.CollectEvents(ev))
	if err != nil {
		return err
	}

	// Check if custom home is provided
	initOptions := initOptionWithHomeFlag(cmd, []networkbuilder.InitOption{})

	// Get keyring backend
	keyringBackend, err := cmd.Flags().GetString(flagKeyringBackend)
	if err != nil {
		return err
	}
	initOptions = append(initOptions, networkbuilder.InitializationKeyringBackend(chaincmd.KeyringBackend(keyringBackend)))

	// init the blockchain.
	blockchain, err := nb.Init(cmd.Context(), chainID, networkbuilder.SourceChainID(), initOptions...)
	if err != nil {
		return err
	}
	defer blockchain.Cleanup()

	// get blockchain's info.
	info, err := blockchain.Info()
	if err != nil {
		return err
	}

	s.Stop()

	// hold default values and user inputs for target chain to later use these to join to the chain.
	var (
		account      *chain.Account
		accountName  = "alice"
		accountCoins = "1000token,100000000stake"
		denom        = "stake"
	)
	if info.Config.Validator.Staked != "" {
		if c, err := types.ParseCoinNormalized(info.Config.Validator.Staked); err == nil {
			denom = c.Denom
		}
	}
	if acc, ok := info.Config.AccountByName(info.Config.Validator.Name); ok {
		accountName = acc.Name
		accountCoins = strings.Join(acc.Coins, ",")
	}

	// ask to propose an account on target blockchain.
	shouldProposeAccount := true

	if gentxPath != "" {
		askAccountProposal := promptui.Prompt{
			Label:     "Would you like to propose an account to place in Genesis",
			IsConfirm: true,
		}
		_, err := askAccountProposal.Run()
		shouldProposeAccount = err == nil
	} else {
		printSection(fmt.Sprintf("Account on the blockchain %s", chainID))

		acc, err := createChainAccount(cmd.Context(), blockchain, fmt.Sprintf("%s blockchain", chainID), accountName)
		if err != nil {
			return err
		}
		account = &acc
	}

	if shouldProposeAccount {
		// ask to create an account proposal.
		printSection("Account proposal")

		if account == nil {
			account = &chain.Account{}
		}

		accQuestions := cliquiz.NewQuestion("Account coins",
			&account.Coins,
			cliquiz.DefaultAnswer(accountCoins),
		)

		if err := cliquiz.Ask(accQuestions); err != nil {
			return err
		}
	}

	gentx, validatorAddress, selfDelegation, publicAddress, err := createValidatorInfo(
		cmd,
		blockchain,
		account,
		denom,
		publicAddress,
		gentxPath,
	)
	if err != nil {
		return err
	}
	if shouldProposeAccount && account.Address == "" {
		account.Address = validatorAddress
	}

	prettyGentx, err := gentx.Pretty()
	if err != nil {
		return err
	}

	// ask to confirm to join to network.
	fmt.Printf("\nYour validator's details (Gentx): \n\n%s\n\n", prettyGentx)
	fmt.Printf("Confirm joining to %s as a validator with the Gentx above:\n", chainID)

	var shouldJoin string
	var answerYes, answerNo = "yes", "no"

	for {
		if err := cliquiz.Ask(cliquiz.NewQuestion(fmt.Sprintf("Please enter %q or %q", answerYes, answerNo),
			&shouldJoin,
			cliquiz.Required(),
		)); err != nil {
			return err
		}
		if shouldJoin == answerNo {
			s.Stop()
			fmt.Println("said no")
			return nil
		}
		if shouldJoin == answerYes {
			break
		}
	}

	// propose to join to the network.
	s.SetText("Proposing...")
	s.Start()

	if err := blockchain.Join(cmd.Context(), account, validatorAddress, publicAddress, gentx, selfDelegation); err != nil {
		return err
	}
	s.Stop()

	fmt.Println("\n📜  Proposal about joining as a validator has been successfully submitted!")
	return nil
}

func createValidatorInfo(
	cmd *cobra.Command,
	blockchain *networkbuilder.Blockchain,
	account *chain.Account,
	denom,
	publicAddress,
	gentxPath string,
) (
	gentx jsondoc.Doc,
	validatorAddress string,
	selfDelegation types.Coin,
	calculatedPublicAddress string,
	err error,
) {
	questions := []cliquiz.Question{}
	var proposal networkbuilder.Proposal

	// prepare questions to interactively ask for validator info when gentx isn't provided.
	if gentxPath == "" {
		questions = append(questions,
			cliquiz.NewQuestion("Staking amount", &proposal.Validator.StakingAmount,
				cliquiz.DefaultAnswer("95000000stake"),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Moniker",
				&proposal.Validator.Moniker,
				cliquiz.DefaultAnswer("mynode"),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Commission rate",
				&proposal.Validator.CommissionRate,
				cliquiz.DefaultAnswer("0.10"),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Commission max rate",
				&proposal.Validator.CommissionMaxRate,
				cliquiz.DefaultAnswer("0.20"),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Commission max change rate",
				&proposal.Validator.CommissionMaxChangeRate,
				cliquiz.DefaultAnswer("0.01"),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Min self delegation",
				&proposal.Validator.MinSelfDelegation,
				cliquiz.DefaultAnswer("1"),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Gas prices",
				&proposal.Validator.GasPrices,
				cliquiz.DefaultAnswer("0.025"+denom),
				cliquiz.Required(),
			),
			cliquiz.NewQuestion("Website", &proposal.Meta.Website),
			cliquiz.NewQuestion("Identity", &proposal.Meta.Identity),
			cliquiz.NewQuestion("Details", &proposal.Meta.Details),
		)
	}

	// prepare questions to interactively ask for a publicAddress when peer isn't provided
	// and not running through chisel proxy.
	if publicAddress == "" && !xchisel.IsEnabled() {
		opts := []cliquiz.Option{
			cliquiz.Required(),
		}
		ip, err := ipify.GetIp()
		if err == nil {
			opts = append(opts, cliquiz.DefaultAnswer(fmt.Sprintf("%s:26656", ip)))
		}
		questions = append(questions, cliquiz.NewQuestion("Peer's address", &publicAddress, opts...))
	}

	// interactively ask validator questions if there is a need to collect extra info.
	if len(questions) > 0 {
		fmt.Println()
		printSection("Validator proposal")

		if err := cliquiz.Ask(questions...); err != nil {
			return nil, "", types.Coin{}, "", err
		}
	}

	// issue gentx and return with validator info when gentx isn't provided manually.
	if gentxPath == "" {
		if gentx, err = blockchain.IssueGentx(cmd.Context(), *account, proposal); err != nil {
			return nil, "", types.Coin{}, "", err
		}

		if selfDelegation, err = types.ParseCoinNormalized(proposal.Validator.StakingAmount); err != nil {
			return nil, "", types.Coin{}, "", err
		}

		return gentx, account.Address, selfDelegation, publicAddress, nil
	}

	// gentx is provided manually so use it and return with validator info.
	if gentx, err = os.ReadFile(gentxPath); err != nil {
		return nil, "", types.Coin{}, "", errors.Wrap(err, "cannot open gentx file")
	}

	info, err := networkbuilder.ParseGentx(gentx)
	if err != nil {
		return nil, "", types.Coin{}, "", err
	}
	return gentx, info.DelegatorAddress, info.SelfDelegation, publicAddress, nil
}

func createChainAccount(ctx context.Context, blockchain *networkbuilder.Blockchain, title, defaultAccountName string) (account chain.Account, err error) {
	var (
		createAccount = "Create a new account"
		importAccount = "Import an account from mnemonic"
	)
	var (
		qs = []*survey.Question{
			{
				Name: "account",
				Prompt: &survey.Select{
					Message: "Choose an account:",
					Options: []string{createAccount, importAccount},
				},
			},
		}
		answers = struct {
			Account string `survey:"account"`
		}{}
	)
	if err := survey.Ask(qs, &answers); err != nil {
		if err == terminal.InterruptErr {
			return account, context.Canceled
		}
		return account, err
	}

	switch answers.Account {
	case createAccount:
		var name string
		if err := cliquiz.Ask(cliquiz.NewQuestion("Account name", &name, cliquiz.DefaultAnswer(defaultAccountName), cliquiz.Required())); err != nil {
			return account, err
		}

		if account, err = blockchain.CreateAccount(ctx, chain.Account{Name: name}); err != nil {
			return account, err
		}

		fmt.Printf("\n%s account has been created successfully!\nAccount address: %s \nMnemonic: %s\n\n",
			title,
			account.Address,
			account.Mnemonic,
		)

	case importAccount:
		var name string
		var mnemonic string
		if err := cliquiz.Ask(
			cliquiz.NewQuestion("Account name", &name, cliquiz.DefaultAnswer(defaultAccountName), cliquiz.Required()),
			cliquiz.NewQuestion("Mnemonic", &mnemonic, cliquiz.Required()),
		); err != nil {
			return account, err
		}

		if account, err = blockchain.CreateAccount(ctx, chain.Account{
			Name:     name,
			Mnemonic: mnemonic,
		}); err != nil {
			return account, err
		}
		fmt.Printf("\n%s account has been imported successfully!\nAccount address: %s\n\n", title, account.Address)
	}
	return account, nil
}
