package networktypes_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

var (
	sampleCoins    = sdk.NewCoins(sdk.NewCoin("bar", sdk.NewInt(1000)), sdk.NewCoin("foo", sdk.NewInt(2000)))
	sampleCoinsStr = sampleCoins.String()
)

func TestToGenesisAccount(t *testing.T) {
	tests := []struct {
		name     string
		fetched  launchtypes.GenesisAccount
		expected networktypes.GenesisAccount
	}{
		{
			name: "genesis account",
			fetched: launchtypes.GenesisAccount{
				Address: "spn123",
				Coins:   sampleCoins,
			},
			expected: networktypes.GenesisAccount{
				Address: "spn123",
				Coins:   sampleCoinsStr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				require.EqualValues(t, tt.expected, networktypes.ToGenesisAccount(tt.fetched))
			})
		})
	}
}

func TestToVestingAccount(t *testing.T) {
	tests := []struct {
		name     string
		fetched  launchtypes.VestingAccount
		expected networktypes.VestingAccount
		isError  bool
	}{
		{
			name: "vesting account",
			fetched: launchtypes.VestingAccount{
				Address:         "spn123",
				StartingBalance: sampleCoins,
				VestingOptions: *launchtypes.NewDelayedVesting(
					sampleCoins,
					1000,
				),
			},
			expected: networktypes.VestingAccount{
				Address:         "spn123",
				StartingBalance: sampleCoinsStr,
				Vesting:         sampleCoinsStr,
				EndTime:         1000,
			},
		},
		{
			name: "unrecognized vesting option",
			fetched: launchtypes.VestingAccount{
				Address:         "spn123",
				StartingBalance: sampleCoins,
				VestingOptions: launchtypes.VestingOptions{
					Options: nil,
				},
			},
			isError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				vestingAcc, err := networktypes.ToVestingAccount(tt.fetched)
				require.EqualValues(t, tt.isError, err != nil)
				require.EqualValues(t, tt.expected, vestingAcc)
			})
		})
	}
}

func TestToGenesisValidator(t *testing.T) {
	tests := []struct {
		name     string
		fetched  launchtypes.GenesisValidator
		expected networktypes.GenesisValidator
	}{
		{
			name: "genesis validator",
			fetched: launchtypes.GenesisValidator{
				GenTx: []byte("abc"),
				Peer:  "abc@0.0.0.0",
			},
			expected: networktypes.GenesisValidator{
				Gentx: []byte("abc"),
				Peer:  "abc@0.0.0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				require.EqualValues(t, tt.expected, networktypes.ToGenesisValidator(tt.fetched))
			})
		})
	}
}
