package starportcmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

var spnAddress string

func NewNetwork() *cobra.Command {
	c := &cobra.Command{
		Use:   "network",
		Short: "Create and start blockchains collaboratively",
		Args:  cobra.ExactArgs(1),
	}

	// configure flags.
	c.PersistentFlags().StringVarP(&spnAddress, "spn-address", "s", "localhost:26657", "An SPN node address")

	// add sub commands.
	c.AddCommand(NewNetworkAccount())
	c.AddCommand(NewNetworkChain())
	c.AddCommand(NewNetworkProposal())
	return c
}

func newNetworkBuilder(options ...networkbuilder.Option) (*networkbuilder.Builder, error) {
	var spnoptions []spn.Option
	// use test keyring backend on Gitpod in order to prevent prompting for keyring
	// password. This happens because Gitpod uses containers.
	//
	// when not on Gitpod, OS keyring backend is used which only asks password once.
	if os.Getenv("GITPOD_WORKSPACE_ID") != "" {
		spnoptions = append(spnoptions, spn.Keyring(keyring.BackendTest))
	}
	spnclient, err := spn.New(spnAddress, spnoptions...)
	if err != nil {
		return nil, err
	}
	return networkbuilder.New(spnclient, options...)
}

// ensureSPNAccount ensures that an SPN account has ben set by interactively asking
// users to create, import or pick an account.
func ensureSPNAccount(b *networkbuilder.Builder) error {
	if _, err := b.AccountInUse(); err == nil {
		return nil
	}

	fmt.Println(`To use Starport Network you need an account.
Please, select an account or create a new one.
	`)

	accounts, err := accountNames(b)
	if err != nil {
		return err
	}
	var (
		createAccount = "Create a new account"
		importAccount = "Import an account from mnemonic"
	)
	list := append(accounts, createAccount, importAccount)
	var (
		qs = []*survey.Question{
			{
				Name: "account",
				Prompt: &survey.Select{
					Message: "Choose an account:",
					Options: list,
				},
			},
		}
		answers = struct {
			Account string `survey:"account"`
		}{}
	)
	err = survey.Ask(qs, &answers)
	if err != nil {
		return err
	}

	var chosenAccountName string

	switch answers.Account {
	case createAccount:
		var name string
		if err := cliquiz.Ask(cliquiz.NewQuestion("Account name", &name)); err != nil {
			return err
		}

		acc, err := b.AccountCreate(name, "")
		if err != nil {
			return err
		}
		fmt.Printf(`Starport Network account has been created successfully!
Account address: %s 
Mnemonic: %s 

`, acc.Address, acc.Mnemonic)
		chosenAccountName = name

	case importAccount:
		var name string
		var mnemonic string
		if err := cliquiz.Ask(
			cliquiz.NewQuestion("Account name", &name),
			cliquiz.NewQuestion("Mnemonic", &mnemonic),
		); err != nil {
			return err
		}

		acc, err := b.AccountCreate(name, mnemonic)
		if err != nil {
			return err
		}
		fmt.Printf(`Starport Network account has been imported successfully!
Account address: %s 

`, acc.Address)
		chosenAccountName = name

	default:
		acc, err := b.AccountGet(answers.Account)
		if err != nil {
			return err
		}
		fmt.Printf(`Starport Network account has been selected.
Account address: %s


`, acc.Address)
		chosenAccountName = answers.Account

	}

	return b.AccountUse(chosenAccountName)
}

// accountNames retrieves a name list of accounts in the OS keyring.
func accountNames(b *networkbuilder.Builder) ([]string, error) {
	var names []string
	accounts, err := b.AccountList()
	if err != nil {
		return nil, err
	}
	for _, account := range accounts {
		names = append(names, account.Name)
	}
	return names, nil
}

func ensureSPNAccountHook(cmd *cobra.Command, args []string) error {
	b, err := newNetworkBuilder()
	if err != nil {
		return err
	}
	err = ensureSPNAccount(b)
	if err == context.Canceled {
		return errors.New("aborted")
	}
	return err
}
