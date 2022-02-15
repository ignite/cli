package networktypes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
)

// MainnetAccount represents the campaign mainnet account of a chain on SPN
type MainnetAccount struct {
	Address string `json:"Address"`
	Shares  string `json:"Shares"`
}

// ToMainnetAccount converts a mainnet account data from SPN and returns a MainnetAccount object
func ToMainnetAccount(acc campaigntypes.MainnetAccount) MainnetAccount {
	launch := MainnetAccount{
		Address: acc.Address,
		Shares:  sdk.Coins(acc.Shares).String(),
	}
	return launch
}

// MainnetVestingAccount represents the campaign mainnet vesting account of a chain on SPN
type MainnetVestingAccount struct {
	Address     string `json:"Address"`
	TotalShares string `json:"TotalShares"`
	Vesting     string `json:"Vesting"`
	EndTime     int64  `json:"EndTime"`
}

// ToMainnetVestingAccount converts a mainnet vesting account data from SPN and returns a MainnetVestingAccount object
func ToMainnetVestingAccount(acc campaigntypes.MainnetVestingAccount) MainnetVestingAccount {
	delaydVesting := acc.VestingOptions.GetDelayedVesting()
	launch := MainnetVestingAccount{
		Address:     acc.Address,
		TotalShares: sdk.Coins(delaydVesting.TotalShares).String(),
		Vesting:     sdk.Coins(delaydVesting.Vesting).String(),
		EndTime:     delaydVesting.EndTime,
	}
	return launch
}
