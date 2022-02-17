package network

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// CampaignUpdateName updates the campaign name
func (n Network) CampaignUpdateName(campaignID uint64, name string) error {
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf(
		"Updating the campaign %d name to %s",
		campaignID,
		name,
	)))

	msg := campaigntypes.NewMsgUpdateCampaignName(n.account.Address(networktypes.SPN), name, campaignID)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var launchRes launchtypes.MsgTriggerLaunchResponse
	if err := res.Decode(&launchRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
		"Now the chain %d name is %s",
		campaignID,
		name,
	)))
	return nil
}

// CampaignUpdateTotalShares updates the campaign total shares
func (n Network) CampaignUpdateTotalShares(campaignID uint64, totalShares campaigntypes.Shares) error {
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf(
		"Updating the campaign %d total shares to %s",
		campaignID,
		totalShares.String(),
	)))

	msg := campaigntypes.NewMsgUpdateTotalShares(n.account.Address(networktypes.SPN), campaignID, totalShares)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var launchRes launchtypes.MsgTriggerLaunchResponse
	if err := res.Decode(&launchRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
		"Now the chain %d total shares is %s",
		campaignID,
		totalShares.String(),
	)))
	return nil
}

// CampaignUpdateTotalSupply updates the campaign total supply
func (n Network) CampaignUpdateTotalSupply(campaignID uint64, totalSupply sdk.Coins) error {
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf(
		"Updating the campaign %d total supply to %s",
		campaignID,
		totalSupply.String(),
	)))

	msg := campaigntypes.NewMsgUpdateTotalSupply(n.account.Address(networktypes.SPN), campaignID, totalSupply)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var launchRes launchtypes.MsgTriggerLaunchResponse
	if err := res.Decode(&launchRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
		"Now the chain %d total supply is %s",
		campaignID,
		totalSupply.String(),
	)))
	return nil
}
