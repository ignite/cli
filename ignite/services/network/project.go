package network

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"

	"github.com/ignite/cli/ignite/pkg/cosmoserror"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

type (
	// Prop updates project proposal.
	Prop func(*updateProp)

	// updateProp represents the update project proposal.
	updateProp struct {
		name        string
		metadata    []byte
		totalSupply sdk.Coins
	}
)

// WithProjectName provides a name proposal to update the project.
func WithProjectName(name string) Prop {
	return func(c *updateProp) {
		c.name = name
	}
}

// WithProjectMetadata provides a meta data proposal to update the project.
func WithProjectMetadata(metadata string) Prop {
	return func(c *updateProp) {
		c.metadata = []byte(metadata)
	}
}

// WithProjectTotalSupply provides a total supply proposal to update the project.
func WithProjectTotalSupply(totalSupply sdk.Coins) Prop {
	return func(c *updateProp) {
		c.totalSupply = totalSupply
	}
}

// Project fetches the project from Network.
func (n Network) Project(ctx context.Context, projectID uint64) (networktypes.Project, error) {
	n.ev.Send("Fetching project information", events.ProgressStart())
	res, err := n.campaignQuery.Campaign(ctx, &campaigntypes.QueryGetCampaignRequest{
		CampaignID: projectID,
	})
	if errors.Is(cosmoserror.Unwrap(err), cosmoserror.ErrNotFound) {
		return networktypes.Project{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.Project{}, err
	}
	return networktypes.ToProject(res.Campaign), nil
}

// Projects fetches the projects from Network.
func (n Network) Projects(ctx context.Context) ([]networktypes.Project, error) {
	var projects []networktypes.Project

	n.ev.Send("Fetching projects information", events.ProgressStart())
	res, err := n.campaignQuery.CampaignAll(ctx, &campaigntypes.QueryAllCampaignRequest{})
	if err != nil {
		return projects, err
	}

	// Parse fetched projects
	for _, project := range res.Campaign {
		projects = append(projects, networktypes.ToProject(project))
	}

	return projects, nil
}

// CreateProject creates a project in Network.
func (n Network) CreateProject(ctx context.Context, name, metadata string, totalSupply sdk.Coins) (uint64, error) {
	n.ev.Send(fmt.Sprintf("Creating project %s", name), events.ProgressStart())
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return 0, err
	}

	msgCreateCampaign := campaigntypes.NewMsgCreateCampaign(
		addr,
		name,
		totalSupply,
		[]byte(metadata),
	)
	res, err := n.cosmos.BroadcastTx(ctx, n.account, msgCreateCampaign)
	if err != nil {
		return 0, err
	}

	var createCampaignRes campaigntypes.MsgCreateCampaignResponse
	if err := res.Decode(&createCampaignRes); err != nil {
		return 0, err
	}

	return createCampaignRes.CampaignID, nil
}

// InitializeMainnet Initialize the mainnet of the project.
func (n Network) InitializeMainnet(
	ctx context.Context,
	projectID uint64,
	sourceURL,
	sourceHash string,
	mainnetChainID string,
) (uint64, error) {
	n.ev.Send("Initializing the mainnet project", events.ProgressStart())
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return 0, err
	}

	msg := campaigntypes.NewMsgInitializeMainnet(
		addr,
		projectID,
		sourceURL,
		sourceHash,
		mainnetChainID,
	)

	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return 0, err
	}

	var initMainnetRes campaigntypes.MsgInitializeMainnetResponse
	if err := res.Decode(&initMainnetRes); err != nil {
		return 0, err
	}

	n.ev.Send(fmt.Sprintf("Project %d initialized on mainnet", projectID), events.ProgressFinish())

	return initMainnetRes.MainnetID, nil
}

// UpdateProject updates the project name or metadata.
func (n Network) UpdateProject(
	ctx context.Context,
	id uint64,
	props ...Prop,
) error {
	// Apply the options provided by the user
	p := updateProp{}
	for _, apply := range props {
		apply(&p)
	}

	n.ev.Send(fmt.Sprintf("Updating the project %d", id), events.ProgressStart())
	account, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	msgs := make([]sdk.Msg, 0)
	if p.name != "" || len(p.metadata) > 0 {
		msgs = append(msgs, campaigntypes.NewMsgEditCampaign(
			account,
			id,
			p.name,
			p.metadata,
		))
	}
	if !p.totalSupply.Empty() {
		msgs = append(msgs, campaigntypes.NewMsgUpdateTotalSupply(
			account,
			id,
			p.totalSupply,
		))
	}

	if _, err := n.cosmos.BroadcastTx(ctx, n.account, msgs...); err != nil {
		return err
	}
	n.ev.Send(fmt.Sprintf("Project %d updated", id), events.ProgressFinish())
	return nil
}
