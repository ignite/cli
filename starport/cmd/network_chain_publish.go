package starportcmd

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagTag     = "tag"
	flagBranch  = "branch"
	flagHash    = "hash"
	flagGenesis = "genesis"
	flagNoCheck = "no-check"
)

// NewNetworkChainPublish returns a new command to publish a new chain to start a new network.
func NewNetworkChainPublish() *cobra.Command {
	c := &cobra.Command{
		Use:   "publish [source-url]",
		Short: "Publush a new chain to start a new network",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainPublishHandler,
	}

	c.Flags().String(flagBranch, "", "Git branch to use for the repo")
	c.Flags().String(flagTag, "", "Git tag to use for the repo")
	c.Flags().String(flagHash, "", "Git hash to use for the repo")
	c.Flags().String(flagGenesis, "", "URL to a custom Genesis")
	c.Flags().Bool(flagNoCheck, false, "Skip verifying chain's integrity")
	c.Flags().String(flagFrom, cosmosaccount.DefaultAccount, "Account name to use for sending transactions to SPN")
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
		noCheck, _    = cmd.Flags().GetBool(flagNoCheck)
	)

	s := clispinner.New()
	defer s.Stop()

	ev := events.NewBus()
	go printEvents(ev, s)

	nb, err := newNetwork(cmd, network.CollectEvents(ev))
	if err != nil {
		return err
	}

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
	s.Stop()

	if err == context.Canceled {
		fmt.Println("aborted")
		return nil
	}
	if err != nil {
		return err
	}
	defer blockchain.Cleanup()

	s.Start()

	// create blockchain.
	var createOptions []network.CreateOption
	if genesisURL != "" {
		createOptions = append(createOptions, network.WithCustomGenesisFromURL(genesisURL))
	}
	if noCheck {
		createOptions = append(createOptions, network.WithNoCheck())
	} else if genesisURL != "" {
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

		if err := blockchain.Init(cmd.Context()); err != nil {
			return err
		}
	}

	s.SetText("Publishing...")

	if err := blockchain.Publish(cmd.Context(), createOptions...); err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("%s Network published\n", clispinner.OK)
	return nil
}
