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
				require.EqualValues(t, tt.expected, networktypes.ToGenesisAccount(tt.fetched))
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
				Address: "spn123",
				VestingOptions: *launchtypes.NewDelayedVesting(
					sampleCoins,
					sampleCoins,
					1000,
				),
			},
			expected: networktypes.VestingAccount{
				Address:      "spn123",
				TotalBalance: sampleCoinsStr,
				Vesting:      sampleCoinsStr,
				EndTime:      1000,
			},
		},
		{
			name: "unrecognized vesting option",
			fetched: launchtypes.VestingAccount{
				Address: "spn123",
				VestingOptions: launchtypes.VestingOptions{
					Options: nil,
				},
			},
			isError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
				vestingAcc, err := networktypes.ToVestingAccount(tt.fetched)
				require.EqualValues(t, tt.isError, err != nil)
				require.EqualValues(t, tt.expected, vestingAcc)
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
				require.EqualValues(t, tt.expected, networktypes.ToGenesisValidator(tt.fetched))
		})
	}
}

func TestGenesisInformation_ApplyRequest(t *testing.T) {
	// used as a template for tests
	genesisInformation := networktypes.NewGenesisInformation(
		[]networktypes.GenesisAccount{
			{
				Address: "spn1g50xher44l9hjuatjdfxgv254jh2wgzfs55yu3",
				Coins: "1000foo",
			},
			{
				Address: "spn1sgphx4vxt63xhvgp9wpewajyxeqt04twfj7gcc",
				Coins: "1000bar",
			},
		},
		[]networktypes.VestingAccount{
			{

			},
		},
		[]networktypes.GenesisValidator{
			{

			},
		},
		)

	launchtypes.NewMsgRequestAddAccount()

	tests := []struct {
		name     string
		gi networktypes.GenesisInformation
		r  launchtypes.Request
		err error
	}{
		{
			name: "genesis account request",
		},
		{
			name: "vesting account request",
		},
		{
			name: "genesis validator request",
		},
		{
			name: "genesis account: existing genesis account",
		},
		{
			name: "genesis account: existing vesting account",
		},
		{
			name: "vesting account: existing genesis account",
		},
		{
			name: "vesting account: existing vesting account",
		},
		{
			name: "existing genesis validator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newGi, err := tt.gi.ApplyRequest(tt.r)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}

			// parse difference following request
			switch rc := tt.r.Content.Content.(type) {
			case *launchtypes.RequestContent_GenesisAccount:
				ga := networktypes.ToGenesisAccount(*rc.GenesisAccount)
				got, ok := newGi.GenesisAccounts[ga.Address]
				require.True(t, ok)
				require.EqualValues(t, ga, got)

			case *launchtypes.RequestContent_VestingAccount:
				va, err := networktypes.ToVestingAccount(*rc.VestingAccount)
				require.NoError(t, err)
				got, ok := newGi.VestingAccounts[va.Address]
				require.True(t, ok)
				require.EqualValues(t, va, got)

			case *launchtypes.RequestContent_AccountRemoval:
				_, ok := newGi.GenesisAccounts[rc.AccountRemoval.Address]
				require.False(t, ok)
				_, ok = newGi.VestingAccounts[rc.AccountRemoval.Address]
				require.False(t, ok)

			case *launchtypes.RequestContent_GenesisValidator:
				gv := networktypes.ToGenesisValidator(*rc.GenesisValidator)
				got, ok := newGi.GenesisValidators[gv.Address]
				require.True(t, ok)
				require.EqualValues(t, gv, got)

			case *launchtypes.RequestContent_ValidatorRemoval:
				_, ok := newGi.GenesisAccounts[rc.ValidatorRemoval.ValAddress]
				require.False(t, ok)
			}
		})
	}
}
