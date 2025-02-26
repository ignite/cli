package ignitecmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

const (
	msgMigration             = "Migrating blockchain config file from v%d to v%d..."
	msgMigrationPrefix       = "Your blockchain config version is v%d and the latest is v%d."
	msgMigrationPrompt       = "Would you like to upgrade your config file to v%d"
	msgMigrationBuf          = "Now ignite supports the `buf.build` (https://buf.build) registry to manage the protobuf dependencies. The embed protoc binary was deprecated and, your blockchain is still using it. Would you like to upgrade and add the `buf.build` config files to `proto/` folder"
	msgMigrationBufProtoDir  = "Ignite proto directory path from the chain config doesn't match the proto directory path from the chain `buf.work.yaml`. Do you want to add the proto path `%[1]v` to the directories list from the buf work file"
	msgMigrationBufProtoDirs = "Chain `buf.work.yaml` file contains directories that don't exist anymore (%[1]v). Do you want to delete them"
	msgMigrationAddTools     = "Some required imports are missing in %s file: %s. Would you like to add them"
	msgMigrationRemoveTools  = "File %s contains deprecated imports: %s. Would you like to remove them"
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
		NewChainLint(),
	)

	return c
}

func preRunHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New()
	defer session.End()

	appPath, err := goModulePath(cmd)
	if err != nil {
		return err
	}

	cfg, cfgPath, err := getChainConfig(cmd)
	if err != nil {
		return err
	}

	if err := configMigrationPreRunHandler(cmd, session, appPath, cfgPath); err != nil {
		return err
	}

	// if err := toolsMigrationPreRunHandler(cmd, session, appPath); err != nil {
	// 	return err
	// }

	return bufMigrationPreRunHandler(cmd, session, appPath, cfg.Build.Proto.Path)
}

func toolsMigrationPreRunHandler(cmd *cobra.Command, session *cliui.Session, appPath string) error {
	session.StartSpinner("Checking missing tools...")

	toolsFilename := filepath.Join(appPath, "go.mod")
	f, _, err := xast.ParseFile(filepath.Join(appPath, "go.mod"))
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
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

	return os.WriteFile(toolsFilename, buf.Bytes(), 0o600)
}

func bufMigrationPreRunHandler(cmd *cobra.Command, session *cliui.Session, appPath, protoDir string) error {
	// check if the buf files exist.
	hasFiles, err := chain.CheckBufFiles(appPath, protoDir)
	if err != nil {
		return err
	}

	if !hasFiles {
		if !getYes(cmd) {
			if err := session.AskConfirm(msgMigrationBuf); err != nil {
				return ErrProtocUnsupported
			}
		}

		runner := xgenny.NewRunner(cmd.Context(), appPath)
		sm, err := chain.BoxBufFiles(runner, appPath, protoDir)
		if err != nil {
			return err
		}

		session.Print("\n🎉 buf.build files added: \n\n")
		session.Printf("%s\n\n", strings.Join(sm.CreatedFiles(), "\n"))
	}

	// check if the buf.work.yaml has the same proto path from the config file.
	hasProtoPath, missingPaths, err := chain.CheckBufProtoDir(appPath, protoDir)
	if err != nil {
		return err
	}

	if !hasProtoPath {
		if !getYes(cmd) {
			if err := session.AskConfirm(fmt.Sprintf(msgMigrationBufProtoDir, protoDir)); err != nil {
				return nil
			}
		}

		if err := chain.AddBufProtoDir(appPath, protoDir); err != nil {
			return err
		}
	}

	if len(missingPaths) > 0 {
		if !getYes(cmd) {
			if err := session.AskConfirm(fmt.Sprintf(msgMigrationBufProtoDirs, strings.Join(missingPaths, ", "))); err != nil {
				return nil
			}
		}

		if err := chain.RemoveBufProtoDirs(appPath, missingPaths...); err != nil {
			return err
		}
	}

	return nil
}

func configMigrationPreRunHandler(cmd *cobra.Command, session *cliui.Session, appPath, cfgPath string) error {
	rawCfg, err := os.ReadFile(cfgPath)
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
					return errors.Errorf("stopping because config version v%d is required to run the command", chainconfig.LatestVersion)
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

		if err := os.WriteFile(cfgPath, buf.Bytes(), 0o600); err != nil {
			return errors.Errorf("config file migration failed: %w", err)
		}
	}
	return nil
}
