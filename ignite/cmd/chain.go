package ignitecmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/chainconfig"
)

// NewChain returns a command that groups sub commands related to compiling, serving
// blockchains and so on.
func NewChain() *cobra.Command {
	c := &cobra.Command{
		Use:     "chain [command]",
		Short:   "Build, initialize and start a blockchain node or perform other actions on the blockchain",
		Long:    `Build, initialize and start a blockchain node or perform other actions on the blockchain.`,
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
	}

	c.AddCommand(addConfigMigrationVerifier(NewChainServe()))
	c.AddCommand(addConfigMigrationVerifier(NewChainBuild()))
	c.AddCommand(addConfigMigrationVerifier(NewChainInit()))
	c.AddCommand(addConfigMigrationVerifier(NewChainFaucet()))
	c.AddCommand(addConfigMigrationVerifier(NewChainSimulate()))

	return c
}

// TODO: Refactor config migration verifier
func addConfigMigrationVerifier(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().AddFlagSet(flagSetConfig())
	cmd.Flags().AddFlagSet(flagSetYes())

	preRunFun := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if preRunFun != nil {
			if err := preRunFun(cmd, args); err != nil {
				return err
			}
		}

		appPath := flagGetPath(cmd)
		configPath := getConfig(cmd)

		var err error
		if configPath == "" {
			configPath, err = chainconfig.LocateDefault(appPath)
			if err != nil {
				return err
			}
		}

		// Check if the version of the Config File is the latest
		currentVersion, latest, err := chainconfig.IsConfigLatest(configPath)
		if err != nil {
			return err
		}

		if !latest {
			if !getYes(cmd) {
				var confirmed bool
				message := fmt.Sprintf("The configuration file of the project is at the version %d. The latest version is %d. Would you like to upgrade your configuration file of your project to the latest version?",
					currentVersion, chainconfig.LatestVersion)
				prompt := &survey.Confirm{
					Message: message,
				}
				if err := survey.AskOne(prompt, &confirmed); err != nil || !confirmed {
					return nil
				}
			}
			// Convert the current Config Yaml to the latest version
			return chainconfig.MigrateLatest(configPath)
		}
		return nil
	}
	return cmd
}
