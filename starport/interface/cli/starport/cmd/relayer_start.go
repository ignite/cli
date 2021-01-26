package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	givenPaths := args

	paths, err := xrelayer.ListPaths(cmd.Context())
	if err != nil {
		return err
	}

	var allPaths []string
	for _, path := range paths {
		allPaths = append(allPaths, path.ID)
	}

	pathsToLink := xstrings.AllOrSomeFilter(allPaths, givenPaths)

	if len(pathsToLink) == 0 {
		fmt.Println("No chains found for linking.")
	}

	var linkedPaths, alreadyLinkedPaths []string

	for _, path := range paths {
		if !xstrings.SliceContains(pathsToLink, path.ID) {
			continue
		}

		if path.IsLinked {
			alreadyLinkedPaths = append(alreadyLinkedPaths, path.ID)
			continue
		}

		fmt.Println()
		printSection(fmt.Sprintf("Starting %s", path.ID))

		linked, _, err := xrelayer.Start(cmd.Context(), path.ID)
		if err != nil {
			fmt.Printf("❌  Couldn't link chains for %q path: %s\n", path.ID, err.Error())

			continue
		}

		linkedPaths = append(linkedPaths, linked...)
	}

	if len(alreadyLinkedPaths) != 0 {
		fmt.Printf("⛓   %d chains already linked.\n", len(alreadyLinkedPaths)*2)
	}

	if len(linkedPaths) != 0 {
		fmt.Printf("\n🔌  Linked %d chains.\n", len(linkedPaths)*2)
	}

	return nil
}
