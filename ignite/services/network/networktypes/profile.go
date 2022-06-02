package networktypes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
)

// Validator represents the Validator profile on SPN
type Validator struct {
	Address           string   `json:"Address"`
	OperatorAddresses []string `json:"OperatorAddresses"`
	Identity          string   `json:"Identity"`
	Website           string   `json:"Website"`
	Details           string   `json:"Details"`
	Moniker           string   `json:"Moniker"`
	SecurityContact   string   `json:"SecurityContact"`
}

func (v Validator) ToProfile(
	campaignID uint64,
	vouchers sdk.Coins,
	shares,
	vestingShares campaigntypes.Shares,
	chainShares,
	chainVestingShares []ChainShare,
) Profile {
	return Profile{
		CampaignID:         campaignID,
		Address:            v.Address,
		Identity:           v.Identity,
		Website:            v.Website,
		Details:            v.Details,
		Moniker:            v.Moniker,
		SecurityContact:    v.SecurityContact,
		Vouchers:           vouchers,
		Shares:             shares,
		VestingShares:      vestingShares,
		ChainShares:        chainShares,
		ChainVestingShares: chainVestingShares,
	}
}

// ToValidator converts a Validator data from SPN and returns a Validator object
func ToValidator(val profiletypes.Validator) Validator {
	return Validator{
		Address:           val.Address,
		OperatorAddresses: val.OperatorAddresses,
		Identity:          val.Description.Identity,
		Website:           val.Description.Website,
		Details:           val.Description.Details,
		Moniker:           val.Description.Moniker,
		SecurityContact:   val.Description.SecurityContact,
	}
}

// Coordinator represents the Coordinator profile on SPN
type Coordinator struct {
	CoordinatorID uint64 `json:"ID"`
	Address       string `json:"Address"`
	Active        bool   `json:"Active"`
	Identity      string `json:"Identity"`
	Website       string `json:"Website"`
	Details       string `json:"Details"`
}

func (c Coordinator) ToProfile(
	campaignID uint64,
	vouchers sdk.Coins,
	shares,
	vestingShares campaigntypes.Shares,
	chainShares,
	chainVestingShares []ChainShare,
) Profile {
	return Profile{
		CampaignID:         campaignID,
		Address:            c.Address,
		Identity:           c.Identity,
		Website:            c.Website,
		Details:            c.Details,
		Vouchers:           vouchers,
		Shares:             shares,
		VestingShares:      vestingShares,
		ChainShares:        chainShares,
		ChainVestingShares: chainVestingShares,
	}
}

// ToCoordinator converts a Coordinator data from SPN and returns a Coordinator object
func ToCoordinator(coord profiletypes.Coordinator) Coordinator {
	return Coordinator{
		CoordinatorID: coord.CoordinatorID,
		Address:       coord.Address,
		Active:        coord.Active,
		Identity:      coord.Description.Identity,
		Website:       coord.Description.Website,
		Details:       coord.Description.Details,
	}
}

type (
	// ChainShare represents the share of a chain on SPN
	ChainShare struct {
		LaunchID uint64    `json:"LaunchID"`
		Shares   sdk.Coins `json:"Shares"`
	}

	// Profile represents the address profile on SPN
	Profile struct {
		Address            string
		CampaignID         uint64 `json:",omitempty"`
		Identity           string `json:",omitempty"`
		Website            string `json:",omitempty"`
		Details            string `json:",omitempty"`
		Moniker            string `json:",omitempty"`
		SecurityContact    string `json:",omitempty"`
		Vouchers           sdk.Coins
		Shares             campaigntypes.Shares
		VestingShares      campaigntypes.Shares
		ChainShares        []ChainShare
		ChainVestingShares []ChainShare
	}

	// ProfileAcc represents the address profile method interface
	ProfileAcc interface {
		ToProfile(
			campaignID uint64,
			vouchers sdk.Coins,
			shares,
			vestingShares campaigntypes.Shares,
			ChainShares,
			ChainVestingShares []ChainShare,
		) Profile
	}
)
