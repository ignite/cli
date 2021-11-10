package starportcmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagTag      = "tag"
	flagBranch   = "branch"
	flagHash     = "hash"
	flagGenesis  = "genesis"
	flagCampaign = "campaign"
	flagNoCheck  = "no-check"
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
		campaign, _   = cmd.Flags().GetUint64(flagCampaign)
		noCheck, _    = cmd.Flags().GetBool(flagNoCheck)
	)
	nb, s, shutdown, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer shutdown()

	// initialize the blockchain
	initOptions := initOptionWithHomeFlag(cmd, []network.InitOption{network.MustNotInitializedBefore()})

	var sourceOption network.SourceOption

	if !xurl.IsLocalPath(source) {
		switch {
		case tag != "":
			sourceOption = network.SourceRemoteTag(source, tag)
		case branch != "":
			sourceOption = network.SourceRemoteBranch(source, branch)
		case hash != "":
			sourceOption = network.SourceRemoteHash(source, hash)
		default:
			sourceOption = network.SourceRemote(source)
		}
	}

	// init the chain.
	blockchain, err := nb.Blockchain(cmd.Context(), sourceOption, initOptions...)
	if err != nil {
		return err
	}

	var createOptions []network.CreateOption

	if genesisURL != "" {
		createOptions = append(createOptions, network.WithCustomGenesisFromURL(genesisURL))
	}
	if campaign != 0 {
		createOptions = append(createOptions, network.WithCampaign(campaign))
	}

	if noCheck {
		createOptions = append(createOptions, network.WithNoCheck())
	} else {
		// perform checks for the chain requires to initialize it and therefore erase the current home if it exists
		// we ask the user for confirmation
		ok, err := blockchain.IsHomeDirExist()
		if err != nil {
			return err
		}

		if ok && !getYes(cmd) {
			home, err := blockchain.Home()
			if err != nil {
				return err
			}
			prompt := promptui.Prompt{
				Label: fmt.Sprintf("Data directory for blockchain already exists: %s. Would you like to overwrite it",
					home,
				),
				IsConfirm: true,
			}

			s.Stop()
			if _, err := prompt.Run(); err != nil {
				fmt.Println("said no")
				return nil
			}
			s.Start()
		}
	}

	launchID, campaignID, err := blockchain.Publish(cmd.Context(), createOptions...)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("%s Network published \n", clispinner.OK)
	fmt.Printf("%s Launch ID: %d \n", clispinner.Bullet, launchID)
	fmt.Printf("%s Campaign ID: %d \n", clispinner.Bullet, campaignID)

	return nil
}
