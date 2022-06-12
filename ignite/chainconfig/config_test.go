package chainconfig

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validators:
  - name: user1
    bonded: "100000000stake"
`

	conf, err := Parse(strings.NewReader(confyml))

	require.NoError(t, err)
	require.Equal(t, []Account{
		{
			Name:  "me",
			Coins: []string{"1000token", "100000000stake"},
		},
		{
			Name:  "you",
			Coins: []string{"5000token"},
		},
	}, conf.Accounts)
	require.Equal(t, "user1", conf.Validators[0].Name)
	require.Equal(t, "100000000stake", conf.Validators[0].Bonded)
}

func TestValidatorParse(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validators:
  - name: user1
    bonded: "100000000stake"
`

	conf, err := Parse(strings.NewReader(confyml))

	require.NoError(t, err)
	require.Equal(t, []Validator{
		{
			Name:   "user1",
			Bonded: "100000000stake",
			App: configmap{
				"grpc": configmap{
					"address": "0.0.0.0:9090",
				},

				"grpc-web": configmap{
					"address": "0.0.0.0:9091",
				},

				"api": configmap{
					"address": "0.0.0.0:1317",
				},
			},

			Config: configmap{
				"rpc": configmap{
					"laddr": "0.0.0.0:26657",
				},

				"p2p": configmap{
					"laddr": "0.0.0.0:26656",
				},

				"pprof_laddr": "0.0.0.0:6060",
			},
		},
	}, conf.Validators)
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
validators:
  - name: user1
    bonded: "100000000stake"
`

	conf, err := Parse(strings.NewReader(confyml))

	require.NoError(t, err)
	require.Equal(t, []Account{
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
	}, conf.Accounts)
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
	require.Equal(t, &ValidationError{"at least one validator is required"}, err)
}

func TestFaucetHost(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validators:
  - name: user1
    bonded: "100000000stake"
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
validators:
  - name: user1
    bonded: "100000000stake"
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
validators:
  - name: user1
    bonded: "100000000stake"
faucet:
  host: "0.0.0.0:4600"
  port: 4700
`
	conf, err = Parse(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, ":4700", FaucetHost(conf))
}
