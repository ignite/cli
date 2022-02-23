package networktypes

import (
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
)

// Campaign represents the campaign of a chain on SPN
type Campaign struct {
	ID                 uint64 `json:"ID"`
	Name               string `json:"Name"`
	CoordinatorID      uint64 `json:"CoordinatorID"`
	MainnetID          uint64 `json:"MainnetID"`
	MainnetInitialized bool   `json:"MainnetInitialized"`
	TotalSupply        string `json:"TotalSupply"`
	AllocatedShares    string `json:"AllocatedShares"`
	DynamicShares      bool   `json:"DynamicShares"`
	TotalShares        string `json:"TotalShares"`
}

// ToCampaign converts a campaign data from SPN and returns a Campaign object
func ToCampaign(campaign campaigntypes.Campaign) Campaign {
	return Campaign{
		ID:                 campaign.CampaignID,
		Name:               campaign.CampaignName,
		CoordinatorID:      campaign.CoordinatorID,
		MainnetID:          campaign.MainnetID,
		MainnetInitialized: campaign.MainnetInitialized,
		TotalSupply:        campaign.TotalSupply.String(),
		AllocatedShares:    campaign.AllocatedShares.String(),
		DynamicShares:      campaign.DynamicShares,
		TotalShares:        campaign.TotalShares.String(),
	}
}
