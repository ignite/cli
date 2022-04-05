package starportcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ignite-hq/cli/starport/pkg/clispinner"
	"github.com/ignite-hq/cli/starport/services/network"
	"github.com/ignite-hq/cli/starport/services/network/networkchain"
)

const (
	flagTag      = "tag"
	flagBranch   = "branch"
	flagHash     = "hash"
	flagGenesis  = "genesis"
	flagCampaign = "campaign"
	flagNoCheck  = "no-check"
	flagChainID  = "chain-id"
)

// NewNetworkChainPublish returns a new command to publish a new chain to start a new network.
func NewNetworkChainPublish() *cobra.Command {
	c := &cobra.Command{
		Use:   "publish [source-url]",
		Short: "Publish a new chain to start a new network",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainPublishHandler,
	}

	c.Flags().String(flagBranch, "", "Git branch to use for the repo")
	c.Flags().String(flagTag, "", "Git tag to use for the repo")
	c.Flags().String(flagHash, "", "Git hash to use for the repo")
	c.Flags().String(flagGenesis, "", "URL to a custom Genesis")
	c.Flags().String(flagChainID, "", "Chain ID to use for this network")
	c.Flags().Uint64(flagCampaign, 0, "Campaign ID to use for this network")
	c.Flags().Bool(flagNoCheck, false, "Skip verifying chain's integrity")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func networkChainPublishHandler(cmd *cobra.Command, args []string) error {
	var (
		source        = args[0]
		tag, _        = cmd.Flags().GetString(flagTag)
		branch, _     = cmd.Flags().GetString(flagBranch)
		hash, _       = cmd.Flags().GetString(flagHash)
		genesisURL, _ = cmd.Flags().GetString(flagGenesis)
		chainID, _    = cmd.Flags().GetString(flagChainID)
		campaign, _   = cmd.Flags().GetUint64(flagCampaign)
		noCheck, _    = cmd.Flags().GetBool(flagNoCheck)
	)

	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

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

	// use custom genesis from url if given.
	if genesisURL != "" {
		initOptions = append(initOptions, networkchain.WithGenesisFromURL(genesisURL))
	}

	// init in a temp dir.
	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(homeDir)

	initOptions = append(initOptions, networkchain.WithHome(homeDir))

	// init the chain.
	c, err := nb.Chain(sourceOption, initOptions...)
	if err != nil {
		return err
	}

	var publishOptions []network.PublishOption

	if genesisURL != "" {
		publishOptions = append(publishOptions, network.WithCustomGenesis(genesisURL))
	}

	if campaign != 0 {
		publishOptions = append(publishOptions, network.WithCampaign(campaign))
	}

	// use custom chain id if given.
	if chainID != "" {
		publishOptions = append(publishOptions, network.WithChainID(chainID))
	}

	if noCheck {
		publishOptions = append(publishOptions, network.WithNoCheck())
	} else if err := c.Init(cmd.Context()); err != nil { // initialize the chain for checking.
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	launchID, campaignID, err := n.Publish(cmd.Context(), c, publishOptions...)
	if err != nil {
		return err
	}

	nb.Spinner.Stop()

	fmt.Printf("%s Network published \n", clispinner.OK)
	fmt.Printf("%s Launch ID: %d \n", clispinner.Bullet, launchID)
	fmt.Printf("%s Campaign ID: %d \n", clispinner.Bullet, campaignID)

	return nil
}
