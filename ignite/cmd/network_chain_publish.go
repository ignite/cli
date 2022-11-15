package ignitecmd

import (
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/spn/pkg/chainid"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/xurl"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networkchain"
)

const (
	flagTag            = "tag"
	flagBranch         = "branch"
	flagHash           = "hash"
	flagGenesisURL     = "genesis-url"
	flagGenesisConfig  = "genesis-config"
	flagCampaign       = "campaign"
	flagShares         = "shares"
	flagNoCheck        = "no-check"
	flagChainID        = "chain-id"
	flagMainnet        = "mainnet"
	flagAccountBalance = "account-balance"
	flagRewardCoins    = "reward.coins"
	flagRewardHeight   = "reward.height"
)

// NewNetworkChainPublish returns a new command to publish a new chain to start a new network.
func NewNetworkChainPublish() *cobra.Command {
	c := &cobra.Command{
		Use:   "publish [source-url]",
		Short: "Publish a new chain to start a new network",
		Long: `To begin the process of launching a blockchain with Ignite, a coordinator needs
to publish the information about a blockchain. The only required bit of
information is the URL of the source code of the blockchain.

The following command publishes the information about an example blockchain:

  ignite network chain publish github.com/ignite/example

This command fetches the source code of the blockchain, compiles the binary,
verifies that a blockchain can be started with the binary, and publishes the
information about the blockchain to Ignite. The command returns an integer number
that acts as an identifier of the chain on Ignite.

By publishing a blockchain on Ignite you become the "coordinator" of this
blockchain. A coordinator is an account that has the authority to approve and
reject validator requests, set parameters of the blockchain and trigger the
launch of the chain.
`,
		Args: cobra.ExactArgs(1),
		RunE: networkChainPublishHandler,
	}

	flagSetClearCache(c)
	c.Flags().String(flagBranch, "", "Git branch to use for the repo")
	c.Flags().String(flagTag, "", "Git tag to use for the repo")
	c.Flags().String(flagHash, "", "Git hash to use for the repo")
	c.Flags().String(flagGenesisURL, "", "URL to a custom Genesis")
	c.Flags().String(flagGenesisConfig, "", "Name of an Ignite config file in the repo for custom Genesis")
	c.Flags().String(flagChainID, "", "Chain ID to use for this network")
	c.Flags().Uint64(flagCampaign, 0, "Campaign ID to use for this network")
	c.Flags().Bool(flagNoCheck, false, "Skip verifying chain's integrity")
	c.Flags().String(flagCampaignMetadata, "", "Add a campaign metadata")
	c.Flags().String(flagCampaignTotalSupply, "", "Add a total of the mainnet of a campaign")
	c.Flags().String(flagShares, "", "Add shares for the campaign")
	c.Flags().Bool(flagMainnet, false, "Initialize a mainnet campaign")
	c.Flags().String(flagAccountBalance, "", "Balance for each approved genesis account for the chain")
	c.Flags().String(flagRewardCoins, "", "Reward coins")
	c.Flags().Int64(flagRewardHeight, 0, "Last reward height")
	c.Flags().String(flagAmount, "", "Amount of coins for account request")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().AddFlagSet(flagSetCheckDependencies())

	return c
}

func networkChainPublishHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	var (
		tag, _                    = cmd.Flags().GetString(flagTag)
		branch, _                 = cmd.Flags().GetString(flagBranch)
		hash, _                   = cmd.Flags().GetString(flagHash)
		genesisURL, _             = cmd.Flags().GetString(flagGenesisURL)
		genesisConfig, _          = cmd.Flags().GetString(flagGenesisConfig)
		chainID, _                = cmd.Flags().GetString(flagChainID)
		campaign, _               = cmd.Flags().GetUint64(flagCampaign)
		noCheck, _                = cmd.Flags().GetBool(flagNoCheck)
		campaignMetadata, _       = cmd.Flags().GetString(flagCampaignMetadata)
		campaignTotalSupplyStr, _ = cmd.Flags().GetString(flagCampaignTotalSupply)
		sharesStr, _              = cmd.Flags().GetString(flagShares)
		isMainnet, _              = cmd.Flags().GetBool(flagMainnet)
		accountBalance, _         = cmd.Flags().GetString(flagAccountBalance)
		rewardCoinsStr, _         = cmd.Flags().GetString(flagRewardCoins)
		rewardDuration, _         = cmd.Flags().GetInt64(flagRewardHeight)
		amount, _                 = cmd.Flags().GetString(flagAmount)
	)

	// parse the amount.
	amountCoins, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return errors.Wrap(err, "error parsing amount")
	}

	accountBalanceCoins, err := sdk.ParseCoinsNormalized(accountBalance)
	if err != nil {
		return errors.Wrap(err, "error parsing account balance")
	}

	source, err := xurl.MightHTTPS(args[0])
	if err != nil {
		return fmt.Errorf("invalid source url format: %w", err)
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	if campaign != 0 && campaignTotalSupplyStr != "" {
		return fmt.Errorf("%s and %s flags cannot be set together", flagCampaign, flagCampaignTotalSupply)
	}
	if isMainnet {
		if campaign == 0 && campaignTotalSupplyStr == "" {
			return fmt.Errorf(
				"%s flag requires one of the %s or %s flags to be set",
				flagMainnet,
				flagCampaign,
				flagCampaignTotalSupply,
			)
		}
		if chainID == "" {
			return fmt.Errorf("%s flag requires the %s flag", flagMainnet, flagChainID)
		}
	}

	if chainID != "" {
		chainName, _, err := chainid.ParseGenesisChainID(chainID)
		if err != nil {
			return errors.Wrapf(err, "invalid chain id: %s", chainID)
		}
		if err := chainid.CheckChainName(chainName); err != nil {
			return errors.Wrapf(err, "invalid chain id name: %s", chainName)
		}
	}

	totalSupply, err := sdk.ParseCoinsNormalized(campaignTotalSupplyStr)
	if err != nil {
		return err
	}

	rewardCoins, err := sdk.ParseCoinsNormalized(rewardCoinsStr)
	if err != nil {
		return err
	}

	if (!rewardCoins.Empty() && rewardDuration == 0) ||
		(rewardCoins.Empty() && rewardDuration > 0) {
		return fmt.Errorf("%s and %s flags must be provided together", flagRewardCoins, flagRewardHeight)
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// use source from chosen target.
	var sourceOption networkchain.SourceOption

	switch {
	case tag != "":
		sourceOption = networkchain.SourceRemoteTag(source, tag)
	case branch != "":
		sourceOption = networkchain.SourceRemoteBranch(source, branch)
	case hash != "":
		sourceOption = networkchain.SourceRemoteHash(source, hash)
	default:
		sourceOption = networkchain.SourceRemote(source)
	}

	var initOptions []networkchain.Option

	// cannot use both genesisURL and genesisConfig
	if genesisURL != "" && genesisConfig != "" {
		return errors.New("cannot use both genesis-url and genesis-config for initial genesis." +
			"Please use only one of the options.")
	}

	// use custom genesis from url if given.
	if genesisURL != "" {
		initOptions = append(initOptions, networkchain.WithGenesisFromURL(genesisURL))
	}

	// use custom genesis config if given
	if genesisConfig != "" {
		initOptions = append(initOptions, networkchain.WithGenesisFromConfig(genesisConfig))
	}

	// init in a temp dir.
	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(homeDir)

	initOptions = append(initOptions, networkchain.WithHome(homeDir))

	// prepare publish options
	publishOptions := []network.PublishOption{network.WithMetadata(campaignMetadata)}

	switch {
	case genesisURL != "":
		publishOptions = append(publishOptions, network.WithCustomGenesisURL(genesisURL))
	case genesisConfig != "":
		publishOptions = append(publishOptions, network.WithCustomGenesisConfig(genesisConfig))

	}

	if campaign != 0 {
		publishOptions = append(publishOptions, network.WithCampaign(campaign))
	} else if campaignTotalSupplyStr != "" {
		totalSupply, err := sdk.ParseCoinsNormalized(campaignTotalSupplyStr)
		if err != nil {
			return err
		}
		if !totalSupply.Empty() {
			publishOptions = append(publishOptions, network.WithTotalSupply(totalSupply))
		}
	}

	// use custom chain id if given.
	if chainID != "" {
		publishOptions = append(publishOptions, network.WithChainID(chainID))
	}

	if !accountBalanceCoins.IsZero() {
		publishOptions = append(publishOptions, network.WithAccountBalance(accountBalanceCoins))
	}

	if isMainnet {
		publishOptions = append(publishOptions, network.Mainnet())
	}

	if !totalSupply.Empty() {
		publishOptions = append(publishOptions, network.WithTotalSupply(totalSupply))
	}

	if sharesStr != "" {
		sharePercentages, err := network.ParseSharePercents(sharesStr)
		if err != nil {
			return err
		}

		publishOptions = append(publishOptions, network.WithPercentageShares(sharePercentages))
	}

	// TODO: Issue an error or warning when this flag is used with "no-check"?
	//       The "check-dependencies" flag is ignored when the "no-check" one is present.
	if flagGetCheckDependencies(cmd) {
		initOptions = append(initOptions, networkchain.CheckDependencies())
	}

	// init the chain.
	c, err := nb.Chain(sourceOption, initOptions...)
	if err != nil {
		return err
	}

	if !noCheck {
		if err := c.Init(cmd.Context(), cacheStorage); err != nil {
			// initialize the chain for checking.
			return fmt.Errorf("blockchain init failed: %w", err)
		}
	}

	session.StartSpinner("Publishing...")

	n, err := nb.Network()
	if err != nil {
		return err
	}

	launchID, campaignID, err := n.Publish(cmd.Context(), c, publishOptions...)
	if err != nil {
		return err
	}

	if !rewardCoins.IsZero() && rewardDuration > 0 {
		if err := n.SetReward(cmd.Context(), launchID, rewardDuration, rewardCoins); err != nil {
			return err
		}
	}

	if !amountCoins.IsZero() {
		if err := n.SendAccountRequestForCoordinator(cmd.Context(), launchID, amountCoins); err != nil {
			return err
		}
	}

	session.Printf("%s Network published \n", icons.OK)
	if isMainnet {
		session.Printf("%s Mainnet ID: %d \n", icons.Bullet, launchID)
	} else {
		session.Printf("%s Launch ID: %d \n", icons.Bullet, launchID)
	}
	if campaignID != 0 {
		session.Printf("%s Campaign ID: %d \n", icons.Bullet, campaignID)
	}

	return nil
}
