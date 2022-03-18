package networktypes

import (
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
