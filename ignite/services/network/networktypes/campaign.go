package networktypes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
)

// Campaign represents the campaign of a chain on SPN
type Campaign struct {
	ID                 uint64    `json:"ID"`
	Name               string    `json:"Name"`
	CoordinatorID      uint64    `json:"CoordinatorID"`
	MainnetID          uint64    `json:"MainnetID"`
	MainnetInitialized bool      `json:"MainnetInitialized"`
	TotalSupply        sdk.Coins `json:"TotalSupply"`
	AllocatedShares    string    `json:"AllocatedShares"`
	Metadata           string    `json:"Metadata"`
}

// ToCampaign converts a campaign data from SPN and returns a Campaign object
func ToCampaign(campaign campaigntypes.Campaign) Campaign {
	return Campaign{
		ID:                 campaign.CampaignID,
		Name:               campaign.CampaignName,
		CoordinatorID:      campaign.CoordinatorID,
		MainnetID:          campaign.MainnetID,
		MainnetInitialized: campaign.MainnetInitialized,
		TotalSupply:        campaign.TotalSupply,
		AllocatedShares:    campaign.AllocatedShares.String(),
		Metadata:           string(campaign.Metadata),
	}
}

// MainnetAccount represents the campaign mainnet account of a chain on SPN
type MainnetAccount struct {
	Address string               `json:"Address"`
	Shares  campaigntypes.Shares `json:"Shares"`
}

// ToMainnetAccount converts a mainnet account data from SPN and returns a MainnetAccount object
func ToMainnetAccount(acc campaigntypes.MainnetAccount) MainnetAccount {
	return MainnetAccount{
		Address: acc.Address,
		Shares:  acc.Shares,
	}
}

// CampaignChains represents the chains of a campaign on SPN
type CampaignChains struct {
	CampaignID uint64   `json:"CampaignID"`
	Chains     []uint64 `json:"Chains"`
}

// ToCampaignChains converts a campaign chains data from SPN and returns a CampaignChains object
func ToCampaignChains(c campaigntypes.CampaignChains) CampaignChains {
	return CampaignChains{
		CampaignID: c.CampaignID,
		Chains:     c.Chains,
	}
}
