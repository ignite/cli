package chainconfig

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig/common"
	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
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
		&v1.Validator{
			Name:   "user1",
			Bonded: "100000000stake",
			App: map[string]interface{}{"grpc": map[string]interface{}{"address": "0.0.0.0:9090"},
				"grpc-web": map[string]interface{}{"address": "0.0.0.0:9091"}, "api": map[string]interface{}{"address": "0.0.0.0:1317"}},
			Config: map[string]interface{}{"rpc": map[string]interface{}{"laddr": "0.0.0.0:26657"},
				"p2p": map[string]interface{}{"laddr": "0.0.0.0:26656"}, "pprof_laddr": "0.0.0.0:6060"},
		}}, conf.ListValidators())

	require.Equal(t, common.Host{
		RPC:     "0.0.0.0:26657",
		P2P:     "0.0.0.0:26656",
		Prof:    "0.0.0.0:6060",
		GRPC:    "0.0.0.0:9090",
		GRPCWeb: "0.0.0.0:9091",
		API:     "0.0.0.0:1317",
	}, conf.GetHost())
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
		&v1.Validator{
			Name:   "user1",
			Bonded: "100000000stake",
			App: map[string]interface{}{"grpc": map[string]interface{}{"address": "0.0.0.0:9090"},
				"grpc-web": map[string]interface{}{"address": "0.0.0.0:9091"}, "api": map[string]interface{}{"address": "0.0.0.0:1317"}},
			Config: map[string]interface{}{"rpc": map[string]interface{}{"laddr": "0.0.0.0:26657"},
				"p2p": map[string]interface{}{"laddr": "0.0.0.0:26656"}, "pprof_laddr": "0.0.0.0:6060"},
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
	require.Equal(t, "0.0.0.0:4600", common.FaucetHost(conf))

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
	require.Equal(t, ":4700", common.FaucetHost(conf))

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
	require.Equal(t, ":4700", common.FaucetHost(conf))
}

func TestParseWithVersion(t *testing.T) {
	tests := []struct {
		TestName        string
		Input           string
		ExpectedError   error
		ExpectedVersion common.Version
	}{{
		TestName: "Parse the config yaml with the field version 0",
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
	}, {
		TestName: "Parse the config yaml with the field version 1",
		Input: `
version: 1
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validators:
  - name: user1
    staked: "100000000stake"
    app:
      grpc:
        address: localhost:8080
      api:
        address: localhost:80801
`,
		ExpectedError:   nil,
		ExpectedVersion: 1,
	}, {
		TestName: "Parse the config yaml with unsupported version",
		Input: `
version: 10000
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validators:
  - name: user1
    bonded: "100000000stake"
`,
		ExpectedError:   &UnsupportedVersionError{Message: "the version is not available in the supported list"},
		ExpectedVersion: 0,
	}}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			conf, err := Parse(strings.NewReader(test.Input))
			if conf != nil {
				require.Equal(t, test.ExpectedVersion, conf.GetVersion())
			}
			require.Equal(t, test.ExpectedError, err)
		})
	}
}

func TestValidator(t *testing.T) {
	tests := []struct {
		TestName                string
		Input                   string
		ExpectedFirstValidator  *v1.Validator
		ExpectedSecondValidator *v1.Validator
	}{{
		TestName: "Parse the config yaml with no addresses for the validator",
		Input: `
version: 1
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validators:
  - name: user1
    bonded: "100000000stake"
  - name: user2
    bonded: "100000000stake"
`,
		ExpectedFirstValidator: &v1.Validator{
			Name:   "user1",
			Bonded: "100000000stake",
			App: map[string]interface{}{"grpc": map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.GRPCPort)},
				"grpc-web": map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.GRPCWebPort)},
				"api":      map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.APIPort)}},
			Config: map[string]interface{}{"rpc": map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", v1.RPCPort)},
				"p2p":         map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", v1.P2P)},
				"pprof_laddr": fmt.Sprintf("0.0.0.0:%d", v1.PPROFPort)},
		},
		ExpectedSecondValidator: &v1.Validator{
			Name:   "user2",
			Bonded: "100000000stake",
			App: map[string]interface{}{"grpc": map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.GRPCPort+v1.DefaultPortMargin)},
				"grpc-web": map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.GRPCWebPort+v1.DefaultPortMargin)},
				"api":      map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.APIPort+v1.DefaultPortMargin)}},
			Config: map[string]interface{}{"rpc": map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", v1.RPCPort+v1.DefaultPortMargin)},
				"p2p":         map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", v1.P2P+v1.DefaultPortMargin)},
				"pprof_laddr": fmt.Sprintf("0.0.0.0:%d", v1.PPROFPort+v1.DefaultPortMargin)},
		},
	}, {
		TestName: "Parse the config yaml with all the addresses for the validator",
		Input: `
version: 1
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validators:
  - name: user1
    bonded: "100000000stake"
    app:
      grpc:
        address: localhost:8080
      api:
        address: localhost:80801
      grpc-web:
        address: localhost:80802
    config:
      rpc:
        laddr: localhost:80807
      p2p:
        laddr: localhost:80804
      pprof_laddr: localhost:80809
  - name: user2
    bonded: "100000000stake"
    app:
      grpc:
        address: localhost:8180
      api:
        address: localhost:81801
      grpc-web:
        address: localhost:81802
    config:
      rpc:
        laddr: localhost:81807
      p2p:
        laddr: localhost:81804
      pprof_laddr: localhost:81809
`,
		ExpectedFirstValidator: &v1.Validator{
			Name:   "user1",
			Bonded: "100000000stake",
			App: map[string]interface{}{"grpc": map[interface{}]interface{}{"address": "localhost:8080"},
				"grpc-web": map[interface{}]interface{}{"address": "localhost:80802"},
				"api":      map[interface{}]interface{}{"address": "localhost:80801"}},
			Config: map[string]interface{}{"rpc": map[interface{}]interface{}{"laddr": "localhost:80807"},
				"p2p":         map[interface{}]interface{}{"laddr": "localhost:80804"},
				"pprof_laddr": "localhost:80809"},
		},
		ExpectedSecondValidator: &v1.Validator{
			Name:   "user2",
			Bonded: "100000000stake",
			App: map[string]interface{}{"grpc": map[interface{}]interface{}{"address": "localhost:8180"},
				"grpc-web": map[interface{}]interface{}{"address": "localhost:81802"},
				"api":      map[interface{}]interface{}{"address": "localhost:81801"}},
			Config: map[string]interface{}{"rpc": map[interface{}]interface{}{"laddr": "localhost:81807"},
				"p2p":         map[interface{}]interface{}{"laddr": "localhost:81804"},
				"pprof_laddr": "localhost:81809"},
		},
	}}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			conf, err := Parse(strings.NewReader(test.Input))
			require.NoError(t, err)
			require.Equal(t, common.Version(1), conf.GetVersion())
			require.Equal(t, test.ExpectedFirstValidator, conf.ListValidators()[0])
			require.Equal(t, test.ExpectedSecondValidator, conf.ListValidators()[1])
		})
	}
}

func TestParseMapInterface(t *testing.T) {
	confyml := `
version: 1
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
validator:
  name: user1
  staked: "100000000stake"
validators:
  - name: user1
    staked: "100000000stake"
    app:
      grpc:
        address: "localhost:8080"
      api:
        address: "localhost:80801"
faucet:
  host: "0.0.0.0:4600"
  port: 4700
init:
  app:
    test-key: test-val:120
`

	_, err := Parse(strings.NewReader(confyml))
	assert.Nil(t, err)
}
