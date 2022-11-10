package ignitecmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/relayer"
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

	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())

	return c
}

func relayerConnectHandler(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		err = handleRelayerAccountErr(err)
	}()

	session := cliui.New()
	defer session.End()

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
		cosmosaccount.WithHome(getKeyringDir(cmd)),
	)
	if err != nil {
		return err
	}

	if err := ca.EnsureDefaultAccount(); err != nil {
		return err
	}

	var (
		use []string
		ids = args
		r   = relayer.New(ca)
	)

	all, err := r.ListPaths(cmd.Context())
	if err != nil {
		return err
	}

	// if no path ids provided, then we connect all of them otherwise,
	// only connect the specified ones.
	if len(ids) == 0 {
		for _, path := range all {
			use = append(use, path.ID)
		}
	} else {
		for _, id := range ids {
			for _, path := range all {
				if id == path.ID {
					use = append(use, path.ID)
					break
				}
			}
		}
	}

	if len(use) == 0 {
		return session.Println("No chains found to connect.")
	}

	session.StartSpinner("Creating links between chains...")

	if err := r.LinkPaths(cmd.Context(), use...); err != nil {
		return err
	}

	session.StopSpinner()

	if err := printSection(session, "Paths"); err != nil {
		return err
	}

	for _, id := range use {
		session.StartSpinner("Loading...")

		path, err := r.GetPath(cmd.Context(), id)
		if err != nil {
			return err
		}

		session.StopSpinner()

		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "%s:\n", path.ID)
		fmt.Fprintf(w, "   \t%s\t>\t(port: %s)\t(channel: %s)\n", path.Src.ChainID, path.Src.PortID, path.Src.ChannelID)
		fmt.Fprintf(w, "   \t%s\t>\t(port: %s)\t(channel: %s)\n", path.Dst.ChainID, path.Dst.PortID, path.Dst.ChannelID)
		fmt.Fprintln(w)
		w.Flush()
		session.Print(buf.String())
	}

	if err := printSection(session, "Listening and relaying packets between chains..."); err != nil {
		return err
	}

	return r.StartPaths(cmd.Context(), use...)
}
