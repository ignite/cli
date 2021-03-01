package starportcmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/services/chain"
)

const (
	flagRebuildProtoOnce = "rebuild-proto-once"
)

// NewBuild returns a new build command to build a blockchain app.
func NewBuild() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Builds and installs an app and its dependencies",
		Args:  cobra.ExactArgs(0),
		RunE:  buildHandler,
	}
	c.Flags().AddFlagSet(flagSetHomes())
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().Bool(flagRebuildProtoOnce, false, "Only enables proto code generation including 3rd party modules")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func buildHandler(cmd *cobra.Command, args []string) error {
	isRebuildProtoOnce, err := cmd.Flags().GetBool(flagRebuildProtoOnce)
	if err != nil {
		return err
	}

	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	if isRebuildProtoOnce {
		chainOption = append(chainOption, chain.EnableThirdPartyModuleCodegen())
	}

	c, err := newChainWithHomeFlags(cmd, appPath, chainOption...)
	if err != nil {
		return err
	}

	if err := c.Build(cmd.Context()); err != nil {
		return err
	}

	binaries, err := c.Binaries()
	if err != nil {
		return err
	}

	fmt.Printf("🗃  Installed. Use with: %s\n", infoColor(strings.Join(binaries, ", ")))

	return nil
}
