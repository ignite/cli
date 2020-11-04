package starportcmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

func NewNetworkAccountUse() *cobra.Command {
	c := &cobra.Command{
		Use:   "use",
		Short: "Pick an account to be used with Starport Network",
		Long: `Pick one of the accounts in OS keyring to put into use or provide one with --name flag.
Picked account will be used while interacting with Starport Network.`,
		RunE: networkAccountUseHandler,
	}
	c.Flags().StringP("name", "n", "", "Account name to put into use")
	return c
}

func networkAccountUseHandler(cmd *cobra.Command, args []string) error {
	b, err := networkbuilder.New(spnAddress)
	if err != nil {
		return err
	}
	name, _ := cmd.Flags().GetString("name")

	// when name is not provided by the flag,
	// list all accounts for user to pick one.
	if name == "" {
		var names []string
		accounts, err := b.AccountList()
		if err != nil {
			return err
		}
		for _, account := range accounts {
			names = append(names, account.Name)
		}
		var (
			qs = []*survey.Question{
				{
					Name: "account",
					Prompt: &survey.Select{
						Message: "Choose an account:",
						Options: names,
					},
				},
			}
			answers = struct {
				AccountName string `survey:"account"`
			}{}
		)
		err = survey.Ask(qs, &answers)
		if err == terminal.InterruptErr {
			fmt.Println("aborted")
			return nil
		}
		if err != nil {
			return err
		}
		name = answers.AccountName
	}
	if err := b.AccountUse(name); err != nil {
		return err
	}
	fmt.Printf("📒 Account put into use: %s\n", color.New(color.FgYellow).SprintFunc()(name))
	return nil
}
