package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const sdkVersionFlag = "sdk-version"

// NewApp creates new command named `app` to create Cosmos scaffolds customized
// by the user given options.
func NewApp() *cobra.Command {
	c := &cobra.Command{
		Use:   "app [github.com/org/repo]",
		Short: "Generates an empty application",
		Args:  cobra.ExactArgs(1),
		RunE:  appHandler,
	}
	c.Flags().String("address-prefix", "cosmos", "Address prefix")
	addSdkVersionFlag(c)
	return c
}

func appHandler(cmd *cobra.Command, args []string) error {
	name := args[0]
	addressPrefix, _ := cmd.Flags().GetString("address-prefix")
	version, err := sdkVersion(cmd)
	if err != nil {
		return err
	}
	sc := scaffolder.New("",
		scaffolder.AddressPrefix(addressPrefix),
		scaffolder.SdkVersion(version),
	)
	path, err := sc.Init(name)
	if err != nil {
		return err
	}
	message := `
⭐️ Successfully created a Cosmos app '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% starport serve

NOTE: add --verbose flag for verbose (detailed) output.
`
	fmt.Printf(message, path)
	return nil
}

func addSdkVersionFlag(c *cobra.Command) {
	c.Flags().String(sdkVersionFlag, string(cosmosver.Stargate), fmt.Sprintf("Target Cosmos-SDK Version %s", cosmosver.MajorVersions))
}

func sdkVersion(c *cobra.Command) (cosmosver.MajorVersion, error) {
	v, _ := c.Flags().GetString(sdkVersionFlag)
	parsed, err := cosmosver.MajorVersions.Parse(v)
	if err != nil {
		return "", fmt.Errorf("%q is an unkown sdk version", v)
	}
	return parsed, nil
}
