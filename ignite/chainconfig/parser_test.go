package chainconfig

import (
	"strings"
	"testing"

	"github.com/ignite/cli/ignite/chainconfig/common"
	v0 "github.com/ignite/cli/ignite/chainconfig/v0"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validator:
  name: user1
  staked: "100000000stake"
`

	conf, err := Parse(strings.NewReader(confyml))

	require.NoError(t, err)
	require.Equal(t, []common.Account{
		{
			Name:  "me",
			Coins: []string{"1000token", "100000000stake"},
		},
		{
			Name:  "you",
			Coins: []string{"5000token"},
		},
	}, conf.ListAccounts())
	require.Equal(t, []common.Validator{
		&v0.Validator{
			Name:   "user1",
			Staked: "100000000stake",
		}}, conf.ListValidators())
}

func TestCoinTypeParse(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
    mnemonic: ozone unfold device pave lemon potato omit insect column wise cover hint narrow large provide kidney episode clay notable milk mention dizzy muffin crazy
    cointype: 7777777
  - name: you
    coins: ["5000token"]
    cointype: 123456
validator:
  name: user1
  staked: "100000000stake"
`

	conf, err := Parse(strings.NewReader(confyml))

	require.NoError(t, err)
	require.Equal(t, []common.Account{
		{
			Name:     "me",
			Coins:    []string{"1000token", "100000000stake"},
			Mnemonic: "ozone unfold device pave lemon potato omit insect column wise cover hint narrow large provide kidney episode clay notable milk mention dizzy muffin crazy",
			CoinType: "7777777",
		},
		{
			Name:     "you",
			Coins:    []string{"5000token"},
			CoinType: "123456",
		},
	}, conf.ListAccounts())
	require.Equal(t, []common.Validator{
		&v0.Validator{
			Name:   "user1",
			Staked: "100000000stake",
		}}, conf.ListValidators())
}

func TestParseInvalid(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
`

	_, err := Parse(strings.NewReader(confyml))
	require.Equal(t, &ValidationError{"validator is required"}, err)
}

func TestFaucetHost(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validator:
  name: user1
  staked: "100000000stake"
faucet:
  host: "0.0.0.0:4600"
`
	conf, err := Parse(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, "0.0.0.0:4600", FaucetHost(conf))

	confyml = `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validator:
  name: user1
  staked: "100000000stake"
faucet:
  port: 4700
`
	conf, err = Parse(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, ":4700", FaucetHost(conf))

	// Port must be higher priority
	confyml = `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validator:
  name: user1
  staked: "100000000stake"
faucet:
  host: "0.0.0.0:4600"
  port: 4700
`
	conf, err = Parse(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, ":4700", FaucetHost(conf))
}

func TestParseWithVersion(t *testing.T) {
	tests := []struct {
		TestName        string
		Input           string
		ExpectedError   error
		ExpectedVersion int
	}{{
		TestName: "Parse the config yaml with the field version",
		Input: `
version: 0
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validator:
  name: user1
  staked: "100000000stake"
`,
		ExpectedError:   nil,
		ExpectedVersion: 1,
	}}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			conf, err := Parse(strings.NewReader(test.Input))
			require.Equal(t, test.ExpectedError, err)
			require.Equal(t, test.ExpectedVersion, conf.GetVersion())
		})
	}
}
