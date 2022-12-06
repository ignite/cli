package networktypes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
)

// Project represents the project of a chain on SPN
type Project struct {
	ID                 uint64    `json:"ID"`
	Name               string    `json:"Name"`
	CoordinatorID      uint64    `json:"CoordinatorID"`
	MainnetID          uint64    `json:"MainnetID"`
	MainnetInitialized bool      `json:"MainnetInitialized"`
	TotalSupply        sdk.Coins `json:"TotalSupply"`
	AllocatedShares    string    `json:"AllocatedShares"`
	Metadata           string    `json:"Metadata"`
}

// ToProject converts a project data from SPN and returns a Project object
func ToProject(campaign campaigntypes.Campaign) Project {
	return Project{
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

// MainnetAccount represents the project mainnet account of a chain on SPN
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

// ProjectChains represents the chains of a project on SPN
type ProjectChains struct {
	ProjectID uint64   `json:"ProjectID"`
	Chains    []uint64 `json:"Chains"`
}

// ToProjectChains converts a project chains data from SPN and returns a ProjectChains object
func ToProjectChains(c campaigntypes.CampaignChains) ProjectChains {
	return ProjectChains{
		ProjectID: c.CampaignID,
		Chains:    c.Chains,
	}
}
