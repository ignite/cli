package ignitecmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/ignite/pkg/goanalysis"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/xast"
	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/doctor"
)

const (
	msgMigration            = "Migrating blockchain config file from v%d to v%d..."
	msgMigrationPrefix      = "Your blockchain config version is v%d and the latest is v%d."
	msgMigrationPrompt      = "Would you like to upgrade your config file to v%d"
	msgMigrationBuf         = "Now ignite supports the `buf.build` (https://buf.build) registry to manage the protobuf dependencies. The embed protoc binary was deprecated and, your blockchain is still using it. Would you like to upgrade and add the `buf.build` config files to `proto/` folder"
	msgMigrationAddTools    = "Some required imports are missing in %s file: %s. Would you like to add them"
	msgMigrationRemoveTools = "File %s contains deprecated imports: %s. Would you like to remove them"
)

var ErrProtocUnsupported = errors.New("code generation using protoc is only supported by Ignite CLI v0.26.1 or older")

// NewChain returns a command that groups sub commands related to compiling, serving
// blockchains and so on.
func NewChain() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain [command]",
		Short: "Build, init and start a blockchain node",
		Long: `Commands in this namespace let you to build, initialize, and start your
blockchain node locally for development purposes.

To run these commands you should be inside the project's directory so that
Ignite can find the source code. To ensure that you are, run "ls", you should
see the following files in the output: "go.mod", "x", "proto", "app", etc.

By default the "build" command will identify the "main" package of the project,
install dependencies if necessary, set build flags, compile the project into a
binary and install the binary. The "build" command is useful if you just want
the compiled binary, for example, to initialize and start the chain manually. It
can also be used to release your chain's binaries automatically as part of
continuous integration workflow.

The "init" command will build the chain's binary and use it to initialize a
local validator node. By default the validator node will be initialized in your
$HOME directory in a hidden directory that matches the name of your project.
This directory is called a data directory and contains a chain's genesis file
and a validator key. This command is useful if you want to quickly build and
initialize the data directory and use the chain's binary to manually start the
blockchain. The "init" command is meant only for development purposes, not
production.

The "serve" command builds, initializes, and starts your blockchain locally with
a single validator node for development purposes. "serve" also watches the
source code directory for file changes and intelligently
re-builds/initializes/starts the chain, essentially providing "code-reloading".
The "serve" command is meant only for development purposes, not production.

To distinguish between production and development consider the following.

In production, blockchains often run the same software on many validator nodes
that are run by different people and entities. To launch a blockchain in
production, the validator entities coordinate the launch process to start their
nodes simultaneously.

During development, a blockchain can be started locally on a single validator
node. This convenient process lets you restart a chain quickly and iterate
faster. Starting a chain on a single node in development is similar to starting
a traditional web application on a local server.

The "faucet" command lets you send tokens to an address from the "faucet"
account defined in "config.yml". Alternatively, you can use the chain's binary
to send token from any other account that exists on chain.

The "simulate" command helps you start a simulation testing process for your
chain.
`,
		Aliases:           []string{"c"},
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: preRunHandler,
	}

	// Add flags required for the configMigrationPreRunHandler
	c.PersistentFlags().AddFlagSet(flagSetConfig())
	c.PersistentFlags().AddFlagSet(flagSetYes())

	c.AddCommand(
		NewChainServe(),
		NewChainBuild(),
		NewChainInit(),
		NewChainFaucet(),
		NewChainSimulate(),
		NewChainDebug(),
	)

	return c
}

func preRunHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New()
	defer session.End()

	path := flagGetPath(cmd)
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	_, appPath, err := gomodulepath.Find(path)
	if err != nil {
		return err
	}

	if err := configMigrationPreRunHandler(cmd, session, appPath); err != nil {
		return err
	}

	if err := toolsMigrationPreRunHandler(cmd, session, appPath); err != nil {
		return err
	}

	return bufMigrationPreRunHandler(cmd, session, appPath)
}

func toolsMigrationPreRunHandler(cmd *cobra.Command, session *cliui.Session, appPath string) error {
	session.StartSpinner("Checking missing tools...")

	toolsFilename := filepath.Join(appPath, doctor.ToolsFile)
	if _, err := os.Stat(toolsFilename); os.IsNotExist(err) {
		return errors.New("the dependency tools file is missing, run `ignite doctor` and try again")
	}

	f, _, err := xast.ParseFile(toolsFilename)
	if err != nil {
		return err
	}

	missing := cosmosgen.MissingTools(f)
	unused := cosmosgen.UnusedTools(f)

	session.StopSpinner()
	if !getYes(cmd) {
		if len(missing) > 0 {
			question := fmt.Sprintf(
				msgMigrationAddTools,
				toolsFilename,
				strings.Join(missing, ", "),
			)
			if err := session.AskConfirm(question); err != nil {
				missing = []string{}
			}
		}

		if len(unused) > 0 {
			question := fmt.Sprintf(
				msgMigrationRemoveTools,
				toolsFilename,
				strings.Join(unused, ", "),
			)
			if err := session.AskConfirm(question); err != nil {
				unused = []string{}
			}
		}
	}
	if len(missing) == 0 && len(unused) == 0 {
		return nil
	}
	session.StartSpinner("Migrating tools...")

	var buf bytes.Buffer
	if err := goanalysis.UpdateInitImports(f, &buf, missing, unused); err != nil {
		return err
	}

	return os.WriteFile(toolsFilename, buf.Bytes(), 0o644)
}

func bufMigrationPreRunHandler(cmd *cobra.Command, session *cliui.Session, appPath string) error {
	hasFiles := chain.CheckBufFiles(appPath)
	if hasFiles {
		return nil
	}

	if !getYes(cmd) {
		if err := session.AskConfirm(msgMigrationBuf); err != nil {
			return ErrProtocUnsupported
		}
	}

	sm, err := chain.BoxBufFiles(appPath)
	if err != nil {
		return err
	}

	session.Print("\n🎉 buf.build files added: \n\n")
	session.Printf("%s\n\n", strings.Join(sm.CreatedFiles(), "\n"))

	return nil
}

func configMigrationPreRunHandler(cmd *cobra.Command, session *cliui.Session, appPath string) (err error) {
	configPath := getConfig(cmd)
	if configPath == "" {
		if configPath, err = chainconfig.LocateDefault(appPath); err != nil {
			return err
		}
	}

	rawCfg, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	version, err := chainconfig.ReadConfigVersion(bytes.NewReader(rawCfg))
	if err != nil {
		return err
	}

	// Config files with older versions must be migrated to the latest before executing the command
	if version != chainconfig.LatestVersion {
		if !getYes(cmd) {
			prefix := fmt.Sprintf(msgMigrationPrefix, version, chainconfig.LatestVersion)
			question := fmt.Sprintf(msgMigrationPrompt, chainconfig.LatestVersion)

			// Confirm before overwriting the config file
			session.Println(prefix)
			if err := session.AskConfirm(question); err != nil {
				if errors.Is(err, promptui.ErrAbort) {
					return fmt.Errorf("stopping because config version v%d is required to run the command", chainconfig.LatestVersion)
				}

				return err
			}

			// Confirm before migrating the config if there are uncommitted changes
			if err := confirmWhenUncommittedChanges(session, appPath); err != nil {
				return err
			}
		} else {
			session.Printf("%s %s\n", icons.Info, colors.Infof(msgMigration, version, chainconfig.LatestVersion))
		}

		// Convert the current config to the latest version and update the YAML file
		var buf bytes.Buffer
		if err := chainconfig.MigrateLatest(bytes.NewReader(rawCfg), &buf); err != nil {
			return err
		}

		if err := os.WriteFile(configPath, buf.Bytes(), 0o755); err != nil {
			return fmt.Errorf("config file migration failed: %w", err)
		}
	}

	return nil
}
