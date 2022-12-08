package ignitecmd

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/yaml"
	"github.com/ignite/cli/ignite/services/network"
)

const (
	flagProjectName        = "name"
	flagProjectMetadata    = "metadata"
	flagProjectTotalSupply = "total-supply"
)

func NewNetworkProjectUpdate() *cobra.Command {
	c := &cobra.Command{
		Use:   "update [project-id]",
		Short: "Update details fo the project of the project",
		Args:  cobra.ExactArgs(1),
		RunE:  networkProjectUpdateHandler,
	}
	c.Flags().String(flagProjectName, "", "update the project name")
	c.Flags().String(flagProjectMetadata, "", "update the project metadata")
	c.Flags().String(flagProjectTotalSupply, "", "update the total of the mainnet of a project")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkProjectUpdateHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	var (
		projectName, _        = cmd.Flags().GetString(flagProjectName)
		metadata, _           = cmd.Flags().GetString(flagProjectMetadata)
		projectTotalSupply, _ = cmd.Flags().GetString(flagProjectTotalSupply)
	)

	totalSupply, err := sdk.ParseCoinsNormalized(projectTotalSupply)
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse project ID
	projectID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	if projectName == "" && metadata == "" && totalSupply.Empty() {
		return fmt.Errorf("at least one of the flags %s must be provided",
			strings.Join([]string{
				flagProjectName,
				flagProjectMetadata,
				flagProjectTotalSupply,
			}, ", "),
		)
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	var proposals []network.Prop

	if projectName != "" {
		proposals = append(proposals, network.WithProjectName(projectName))
	}
	if metadata != "" {
		proposals = append(proposals, network.WithProjectMetadata(metadata))
	}
	if !totalSupply.Empty() {
		proposals = append(proposals, network.WithProjectTotalSupply(totalSupply))
	}

	if err = n.UpdateProject(cmd.Context(), projectID, proposals...); err != nil {
		return err
	}

	project, err := n.Project(cmd.Context(), projectID)
	if err != nil {
		return err
	}
	session.Println()

	info, err := yaml.Marshal(cmd.Context(), project)
	if err != nil {
		return err
	}

	return session.Print(info)
}
