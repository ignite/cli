package networktypes

import (
	"errors"

	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// GenesisInformation represents all information for a chain to construct the genesis.
type GenesisInformation struct {
	GenesisAccounts   []GenesisAccount
	VestingAccounts   []VestingAccount
	GenesisValidators []GenesisValidator
}

// GenesisAccount represents an account with initial coin allocation for the chain for the chain genesis
type GenesisAccount struct {
	Address string
	Coins   string
}

// VestingAccount represents a vesting account with initial coin allocation  and vesting option for the chain genesis
// VestingAccount supports currently only delayed vesting option
type VestingAccount struct {
	Address      string
	TotalBalance string
	Vesting      string
	EndTime      int64
}

// GenesisValidator represents a genesis validator associated with a gentx in the chain genesis
type GenesisValidator struct {
	Gentx []byte
	Peer  string
}

// NewGenesisInformation initializes a new GenesisInformation
func NewGenesisInformation(
	genAccs []GenesisAccount,
	vestingAccs []VestingAccount,
	genVals []GenesisValidator,
) GenesisInformation {
	return GenesisInformation{
		GenesisAccounts:   genAccs,
		VestingAccounts:   vestingAccs,
		GenesisValidators: genVals,
	}
}

// ToGenesisAccount converts genesis account from SPN
func ToGenesisAccount(acc launchtypes.GenesisAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins.String(),
	}
}

// ToVestingAccount converts vesting account from SPN
func ToVestingAccount(acc launchtypes.VestingAccount) (VestingAccount, error) {
	delayedVesting := acc.VestingOptions.GetDelayedVesting()
	if delayedVesting == nil {
		return VestingAccount{}, errors.New("only delayed vesting option is supported")
	}

	return VestingAccount{
		Address:      acc.Address,
		TotalBalance: delayedVesting.TotalBalance.String(),
		Vesting:      delayedVesting.Vesting.String(),
		EndTime:      delayedVesting.EndTime,
	}, nil
}

// ToGenesisValidator converts genesis validator from SPN
func ToGenesisValidator(val launchtypes.GenesisValidator) GenesisValidator {
	return GenesisValidator{
		Gentx: val.GenTx,
		Peer:  val.Peer,
	}
}
