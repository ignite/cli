package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

// NewChainLint returns a lint command to build a blockchain app.
func NewChainLint() *cobra.Command {
	c := &cobra.Command{
		Use:   "lint",
		Short: "Lint codebase using golangci-lint",
		Long:  "The lint command runs the golangci-lint tool to lint the codebase.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			session := cliui.New(
				cliui.StartSpinnerWithText("Linting..."),
				cliui.WithoutUserInteraction(getYes(cmd)),
			)
			defer session.End()

			chainOption := []chain.Option{
				chain.WithOutputer(session),
				chain.CollectEvents(session.EventBus()),
			}

			c, err := chain.NewWithHomeFlags(cmd, chainOption...)
			if err != nil {
				return err
			}

			return c.Lint(cmd.Context())
		},
	}

	return c
}
