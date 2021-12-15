package networktypes_test

import (
	"testing"
	"time"

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
	newCoin := func(str string) sdk.Coin {
		c, err := sdk.ParseCoinNormalized(str)
		require.NoError(t, err)
		return c
	}
	newCoins := func(str string) sdk.Coins {
		c, err := sdk.ParseCoinsNormalized(str)
		require.NoError(t, err)
		return c
	}

	// used as a template for tests
	genesisInformation := networktypes.NewGenesisInformation(
		[]networktypes.GenesisAccount{
			{
				Address: "spn1g50xher44l9hjuatjdfxgv254jh2wgzfs55yu3",
				Coins:   "1000foo",
			},
		},
		[]networktypes.VestingAccount{
			{
				Address:      "spn1gkzf4e0x6wr4djfd8h82v6cy507gy5v4spaus3",
				TotalBalance: "1000foo",
				Vesting:      "500foo",
				EndTime:      time.Now().Unix(),
			},
		},
		[]networktypes.GenesisValidator{
			{
				Address: "spn1pquxnnpnjyl3ptz3uxs0lrs93s5ljepzq4wyp6",
				Gentx:   []byte("aaa"),
				Peer:    "foo",
			},
		},
	)

	tests := []struct {
		name string
		gi   networktypes.GenesisInformation
		r    launchtypes.Request
		err  error
	}{
		{
			name: "genesis account request",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewGenesisAccount(
					0,
					"spn1sgphx4vxt63xhvgp9wpewajyxeqt04twfj7gcc",
					newCoins("1000bar"),
				),
			},
		},
		{
			name: "vesting account request",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewVestingAccount(
					0,
					"spn19klee4szqpeu0laqze5srhdxtp6fuhcztdrh7c",
					*launchtypes.NewDelayedVesting(
						newCoins("1000bar"),
						newCoins("500bar"),
						time.Now().Unix(),
					),
				),
			},
		},
		{
			name: "genesis validator request",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewGenesisValidator(
					0,
					"spn1xnn9w76mf42t249486ss65lvga7gqs02erpw24",
					[]byte("bbb"),
					[]byte("ccc"),
					newCoin("1000bar"),
					"bar",
				),
			},
		},
		{
			name: "genesis account: existing genesis account",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewGenesisAccount(
					0,
					"spn1g50xher44l9hjuatjdfxgv254jh2wgzfs55yu3",
					newCoins("1000bar"),
				),
			},
			err: networktypes.ErrInvalidRequest,
		},
		{
			name: "genesis account: existing vesting account",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewGenesisAccount(
					0,
					"spn1gkzf4e0x6wr4djfd8h82v6cy507gy5v4spaus3",
					newCoins("1000bar"),
				),
			},
			err: networktypes.ErrInvalidRequest,
		},
		{
			name: "vesting account: existing genesis account",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewVestingAccount(
					0,
					"spn1g50xher44l9hjuatjdfxgv254jh2wgzfs55yu3",
					*launchtypes.NewDelayedVesting(
						newCoins("1000bar"),
						newCoins("500bar"),
						time.Now().Unix(),
					),
				),
			},
			err: networktypes.ErrInvalidRequest,
		},
		{
			name: "vesting account: existing vesting account",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewVestingAccount(
					0,
					"spn1gkzf4e0x6wr4djfd8h82v6cy507gy5v4spaus3",
					*launchtypes.NewDelayedVesting(
						newCoins("1000bar"),
						newCoins("500bar"),
						time.Now().Unix(),
					),
				),
			},
			err: networktypes.ErrInvalidRequest,
		},
		{
			name: "existing genesis validator",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewGenesisValidator(
					0,
					"spn1pquxnnpnjyl3ptz3uxs0lrs93s5ljepzq4wyp6",
					[]byte("bbb"),
					[]byte("ccc"),
					newCoin("1000bar"),
					"bar",
				),
			},
			err: networktypes.ErrInvalidRequest,
		},
		{
			name: "remove genesis account",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewAccountRemoval("spn1g50xher44l9hjuatjdfxgv254jh2wgzfs55yu3"),
			},
		},
		{
			name: "remove vesting account",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewAccountRemoval("spn1gkzf4e0x6wr4djfd8h82v6cy507gy5v4spaus3"),
			},
		},
		{
			name: "remove genesis validator",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewValidatorRemoval("spn1pquxnnpnjyl3ptz3uxs0lrs93s5ljepzq4wyp6"),
			},
		},
		{
			name: "remove account: non-existent account",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewAccountRemoval("spn1pquxnnpnjyl3ptz3uxs0lrs93s5ljepzq4wyp6"),
			},
			err: networktypes.ErrInvalidRequest,
		},
		{
			name: "remove account: non-existent genesis validator",
			gi:   genesisInformation,
			r: launchtypes.Request{
				Content: launchtypes.NewValidatorRemoval("spn1g50xher44l9hjuatjdfxgv254jh2wgzfs55yu3"),
			},
			err: networktypes.ErrInvalidRequest,
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
