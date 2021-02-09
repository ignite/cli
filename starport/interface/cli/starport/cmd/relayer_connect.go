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

	linkedPaths, alreadyLinkedPaths, err := xrelayer.Link(cmd.Context(), pathsToUse...)
	if err != nil {
		return err
	}

	s.Stop()

	if len(alreadyLinkedPaths) != 0 {
		fmt.Printf("⛓  %d paths already created to link chains.\n", len(alreadyLinkedPaths))
	}

	if len(linkedPaths) != 0 {
		fmt.Printf("🔌  Linked chains with %d paths.\n", len(linkedPaths))
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

		rpath := path.Path

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "%s:\n", path.ID)
		fmt.Fprintf(w, "   \t%s\t>\t(port: %s)\t(channel: %s)\n", rpath.Src.ChainID, rpath.Src.PortID, rpath.Src.ChannelID)
		fmt.Fprintf(w, "   \t%s\t>\t(port: %s)\t(channel: %s)\n", rpath.Dst.ChainID, rpath.Dst.PortID, rpath.Dst.ChannelID)
		fmt.Fprintln(w)
		w.Flush()
	}

	printSection("Listening and relaying txs between chains...")

	return xrelayer.Start(cmd.Context(), pathsToUse...)
}
