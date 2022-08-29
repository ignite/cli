package ignitecmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/chainconfig"
)

var migrateMsg = `Your blockchain config version is v%[1]d and the latest is v%[2]d. Would you like to upgrade your config file to v%[2]d?`

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

func addConfigMigrationVerifier(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().AddFlagSet(flagSetConfig())
	cmd.Flags().AddFlagSet(flagSetYes())

	preRunFun := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) (err error) {
		if preRunFun != nil {
			if err = preRunFun(cmd, args); err != nil {
				return err
			}
		}

		configPath := getConfig(cmd)
		if configPath == "" {
			appPath := flagGetPath(cmd)

			if configPath, err = chainconfig.LocateDefault(appPath); err != nil {
				return err
			}
		}

		rawCfg, err := ioutil.ReadFile(configPath)
		if err != nil {
			return err
		}

		version, err := chainconfig.ReadConfigVersion(bytes.NewReader(rawCfg))
		if err != nil {
			return err
		}

		// Config files with older versions must be migrated to the latest before executing the command
		if version != chainconfig.LatestVersion {
			// Confirm before overwritting the config file
			if !getYes(cmd) {
				confirmed := false
				prompt := &survey.Confirm{
					Message: fmt.Sprintf(migrateMsg, version, chainconfig.LatestVersion),
				}

				err := survey.AskOne(prompt, &confirmed)
				if err != nil {
					return err
				} else if !confirmed {
					return fmt.Errorf(
						"stopping because config version v%d is required to run the command",
						chainconfig.LatestVersion,
					)
				}
			}

			file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return err
			}

			defer file.Close()

			// Convert the current config to the latest version and update the YAML file
			return chainconfig.MigrateLatest(bytes.NewReader(rawCfg), file)
		}

		return nil
	}

	return cmd
}
