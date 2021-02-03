package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/xrelayer"
	"github.com/tendermint/starport/starport/pkg/xstrings"
)

// NewRelayerStart returns a new relayer start command to link all or some relayer paths.
// if not paths are specified, all paths are linked.
func NewRelayerStart() *cobra.Command {
	c := &cobra.Command{
		Use:  "start [<path>,...]",
		RunE: relayerStartHandler,
	}
	return c
}

func relayerStartHandler(cmd *cobra.Command, args []string) error {
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

		fmt.Println("No chains found to link.")
		return nil
	}

	s.SetText("Starting...")

	linkedPaths, alreadyLinkedPaths, err := xrelayer.Start(cmd.Context(), pathsToUse...)
	if err != nil {
		return err
	}

	s.Stop()

	if len(alreadyLinkedPaths) != 0 {
		fmt.Printf("⛓  %d chains already linked.\n", len(alreadyLinkedPaths)*2)
	}

	if len(linkedPaths) != 0 {
		fmt.Printf("🔌  Linked %d chains.\n", len(linkedPaths)*2)
	}

	return nil
}
