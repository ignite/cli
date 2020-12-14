package starportcmd

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewNetworkAccountUse creates a new account use command to pick a
// default account to access SPN.
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
	nb, err := newNetworkBuilder()
	if err != nil {
		return err
	}

	name, _ := cmd.Flags().GetString("name")

	// when name is not provided by the flag,
	// list all accounts for user to pick one.
	if name == "" {
		accounts, err := accountNames(nb)
		if err != nil {
			return err
		}
		if len(accounts) == 0 {
			return errors.New("no account found. please create one with 'starport network account create'")
		}
		var (
			qs = []*survey.Question{
				{
					Name: "account",
					Prompt: &survey.Select{
						Message: "Choose an account:",
						Options: accounts,
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
	if err := nb.AccountUse(name); err != nil {
		return err
	}
	fmt.Printf("📒 Account put into use: %s\n", color.New(color.FgYellow).SprintFunc()(name))
	return nil
}
