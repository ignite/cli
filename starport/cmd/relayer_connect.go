package starportcmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/xrelayer"
	"github.com/tendermint/starport/starport/pkg/xstrings"
)

// NewRelayerConnect returns a new relayer connect command to link all or some relayer paths and start
// relaying txs in between.
// if not paths are specified, all paths are linked.
func NewRelayerConnect() *cobra.Command {
	c := &cobra.Command{
		Use:   "connect [<path>,...]",
		Short: "Link chains associated with paths and start relaying tx packets in between",
		RunE:  relayerConnectHandler,
	}
	return c
}

func relayerConnectHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New()
	defer s.Stop()

	allPaths, err := xrelayer.ListPaths(cmd.Context())
	if err != nil {
		return err
	}

	var (
		givenPathIDs = args
		allPathIDs   = xstrings.List(len(allPaths), func(i int) string { return allPaths[i].ID })
		pathsToUse   = xstrings.AllOrSomeFilter(allPathIDs, givenPathIDs)
	)

	if len(pathsToUse) == 0 {
		s.Stop()

		fmt.Println("No chains found to connect.")
		return nil
	}

	s.SetText("Linking paths between chains...")

	linkedPaths, alreadyLinkedPaths, failedToLinkPaths, err := xrelayer.Link(cmd.Context(), pathsToUse...)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Println()
	printSection("Linking chains")

	if len(alreadyLinkedPaths) != 0 {
		fmt.Printf("✓ %d paths already created to link chains.\n", len(alreadyLinkedPaths))
		for _, id := range alreadyLinkedPaths {
			fmt.Printf("  - %s\n", id)
		}
		fmt.Println()
	}

	if len(linkedPaths) != 0 {
		fmt.Printf("✓ Linked chains with %d paths.\n", len(linkedPaths))
		for _, id := range linkedPaths {
			fmt.Printf("  - %s\n", id)
		}
		fmt.Println()
	}

	if len(failedToLinkPaths) != 0 {
		fmt.Printf("x Failed to link chains in %d paths.\n", len(failedToLinkPaths))
		for _, failed := range failedToLinkPaths {
			fmt.Printf("  - %s failed with error: %s\n", failed.ID, failed.ErrorMsg)
		}
		fmt.Println()
		fmt.Printf("Continuing with %d paths...\n", len(alreadyLinkedPaths)+len(linkedPaths))
	}

	fmt.Println()
	printSection("Chains by paths")

	for _, id := range append(linkedPaths, alreadyLinkedPaths...) {
		s.SetText("Loading...").Start()

		path, err := xrelayer.GetPath(cmd.Context(), id)
		if err != nil {
			return err
		}

		s.Stop()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "%s:\n", path.ID)
		fmt.Fprintf(w, "   \t%s\t>\t(port: %s)\t(channel: %s)\n", path.Src.ChainID, path.Src.PortID, path.Src.ChannelID)
		fmt.Fprintf(w, "   \t%s\t>\t(port: %s)\t(channel: %s)\n", path.Dst.ChainID, path.Dst.PortID, path.Dst.ChannelID)
		fmt.Fprintln(w)
		w.Flush()
	}

	printSection("Listening and relaying packets between chains...")

	return xrelayer.Start(cmd.Context(), append(linkedPaths, alreadyLinkedPaths...)...)
}
