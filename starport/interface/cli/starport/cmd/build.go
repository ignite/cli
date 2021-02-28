package starportcmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/services/chain"
)

const (
	onlyProtoFlag = "only-proto"
	protoAllFlag  = "proto-all"
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
	c.Flags().Bool(onlyProtoFlag, false, "Only enables proto code generation")
	c.Flags().Bool(protoAllFlag, false, "Enables proto code generation for dependency modules as well")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func buildHandler(cmd *cobra.Command, args []string) error {
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
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
