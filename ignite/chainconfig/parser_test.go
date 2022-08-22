package chainconfig_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
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

	conf, err := chainconfig.Parse(strings.NewReader(confyml))

	require.NoError(t, err)
	require.Equal(t, []config.Account{
		{
			Name:  "me",
			Coins: []string{"1000token", "100000000stake"},
		},
		{
			Name:  "you",
			Coins: []string{"5000token"},
		},
	}, conf.ListAccounts())

	require.Equal(t, []v1.Validator{
		{
			Name:   "user1",
			Bonded: "100000000stake",
			App: map[string]interface{}{
				"grpc":     map[string]interface{}{"address": "0.0.0.0:9090"},
				"grpc-web": map[string]interface{}{"address": "0.0.0.0:9091"},
				"api":      map[string]interface{}{"address": "0.0.0.0:1317"},
			},
			Config: map[string]interface{}{
				"rpc":         map[string]interface{}{"laddr": "0.0.0.0:26657"},
				"p2p":         map[string]interface{}{"laddr": "0.0.0.0:26656"},
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
validator:
  name: user1
  staked: "100000000stake"
`

	conf, err := chainconfig.Parse(strings.NewReader(confyml))

	require.NoError(t, err)
	require.Equal(t, []config.Account{
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
	require.Equal(t, []v1.Validator{
		{
			Name:   "user1",
			Bonded: "100000000stake",
			App: map[string]interface{}{
				"grpc":     map[string]interface{}{"address": "0.0.0.0:9090"},
				"grpc-web": map[string]interface{}{"address": "0.0.0.0:9091"},
				"api":      map[string]interface{}{"address": "0.0.0.0:1317"},
			},
			Config: map[string]interface{}{
				"rpc":         map[string]interface{}{"laddr": "0.0.0.0:26657"},
				"p2p":         map[string]interface{}{"laddr": "0.0.0.0:26656"},
				"pprof_laddr": "0.0.0.0:6060",
			},
		},
	}, conf.Validators)
}

func TestParseInvalid(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
`

	_, err := chainconfig.Parse(strings.NewReader(confyml))
	require.Equal(t, &chainconfig.ValidationError{"validator is required"}, err)
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
	conf, err := chainconfig.Parse(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, "0.0.0.0:4600", chainconfig.FaucetHost(conf))

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
	conf, err = chainconfig.Parse(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, ":4700", chainconfig.FaucetHost(conf))

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
	conf, err = chainconfig.Parse(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, ":4700", chainconfig.FaucetHost(conf))
}

func TestParseWithMigration(t *testing.T) {
	confyml := `
accounts:
  - name: alice
    coins: ["100000000uatom", "100000000000000000000aevmos"]
  - name: bob
    coins: ["5000000000000aevmos"]
validator:
  name: alice
  staked: "100000000000000000000aevmos"
faucet:
  name: bob 
  coins: ["10aevmos"]
build:
  binary: "evmosd"
init:
  home: "$HOME/.evmosd"
  app:
    evm-rpc:
      address: "0.0.0.0:8545"     # change the JSON-RPC address and port
      ws-address: "0.0.0.0:8546"  # change the JSON-RPC websocket address and port
genesis:
  chain_id: "evmosd_9000-1"
  app_state:
    staking:
      params:
        bond_denom: "aevmos"
    mint:
      params:
        mint_denom: "aevmos"
    crisis:
      constant_fee:
        denom: "aevmos"
    gov:
      deposit_params:
        min_deposit:
          - amount: "10000000"
            denom: "aevmos"
    evm:
      params:
        evm_denom: "aevmos"
`
	conf, err := chainconfig.Parse(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, config.Version(1), conf.GetVersion())
	require.Equal(t, []config.Account{
		{
			Name:  "alice",
			Coins: []string{"100000000uatom", "100000000000000000000aevmos"},
		},
		{
			Name:  "bob",
			Coins: []string{"5000000000000aevmos"},
		},
	}, conf.ListAccounts())
	require.Equal(t, "bob", *conf.Faucet.Name)
	require.Equal(t, []string{"10aevmos"}, conf.Faucet.Coins)
	// The default value of Host has been filled in for Faucet.
	require.Equal(t, "0.0.0.0:4500", conf.Faucet.Host)
	// The default values have been filled in for Build.
	require.Equal(t, config.Build{
		Binary: "evmosd",
		Proto: config.Proto{
			Path: "proto",
			ThirdPartyPaths: []string{
				"third_party/proto",
				"proto_vendor",
			},
		},
	}, conf.Build)

	// The validator is filled with the default values for grpc, grpc-web, api, rpc, p2p and pprof_laddr.
	// The init.app and init.home are moved under the validator as well.
	require.Equal(t, []v1.Validator{
		{
			Name:   "alice",
			Bonded: "100000000000000000000aevmos",
			Home:   "$HOME/.evmosd",
			App: map[string]interface{}{
				"grpc":     map[string]interface{}{"address": "0.0.0.0:9090"},
				"grpc-web": map[string]interface{}{"address": "0.0.0.0:9091"},
				"api":      map[string]interface{}{"address": "0.0.0.0:1317"},
				"evm-rpc": map[interface{}]interface{}{
					"address":    "0.0.0.0:8545",
					"ws-address": "0.0.0.0:8546",
				},
			},
			Config: map[string]interface{}{
				"rpc":         map[string]interface{}{"laddr": "0.0.0.0:26657"},
				"p2p":         map[string]interface{}{"laddr": "0.0.0.0:26656"},
				"pprof_laddr": "0.0.0.0:6060",
			},
		},
	}, conf.Validators)

	require.Equal(t, map[string]interface{}{
		"app_state": map[interface{}]interface{}{
			"crisis": map[interface{}]interface{}{
				"constant_fee": map[interface{}]interface{}{"denom": "aevmos"},
			},
			"evm": map[interface{}]interface{}{
				"params": map[interface{}]interface{}{"evm_denom": "aevmos"},
			},
			"gov": map[interface{}]interface{}{
				"deposit_params": map[interface{}]interface{}{
					"min_deposit": []interface{}{
						map[interface{}]interface{}{
							"amount": "10000000",
							"denom":  "aevmos",
						},
					},
				},
			},
			"mint": map[interface{}]interface{}{
				"params": map[interface{}]interface{}{
					"mint_denom": "aevmos",
				},
			},
			"staking": map[interface{}]interface{}{
				"params": map[interface{}]interface{}{
					"bond_denom": "aevmos",
				},
			},
		},
		"chain_id": "evmosd_9000-1",
	},
		conf.Genesis)
}

func TestParseWithVersion(t *testing.T) {
	tests := []struct {
		TestName        string
		Input           string
		ExpectedError   error
		ExpectedVersion config.Version
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
		ExpectedError:   &chainconfig.UnsupportedVersionError{Message: "the version is not available in the supported list"},
		ExpectedVersion: 0,
	}}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			conf, err := chainconfig.Parse(strings.NewReader(test.Input))
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
		ExpectedFirstValidator  v1.Validator
		ExpectedSecondValidator v1.Validator
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
		ExpectedFirstValidator: v1.Validator{
			Name:   "user1",
			Bonded: "100000000stake",
			App: map[string]interface{}{
				"grpc":     map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.GRPCPort)},
				"grpc-web": map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.GRPCWebPort)},
				"api":      map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.APIPort)},
			},
			Config: map[string]interface{}{
				"rpc":         map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", v1.RPCPort)},
				"p2p":         map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", v1.P2PPort)},
				"pprof_laddr": fmt.Sprintf("0.0.0.0:%d", v1.PPROFPort),
			},
		},
		ExpectedSecondValidator: v1.Validator{
			Name:   "user2",
			Bonded: "100000000stake",
			App: map[string]interface{}{
				"grpc":     map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.GRPCPort+v1.DefaultPortMargin)},
				"grpc-web": map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.GRPCWebPort+v1.DefaultPortMargin)},
				"api":      map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", v1.APIPort+v1.DefaultPortMargin)},
			},
			Config: map[string]interface{}{
				"rpc":         map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", v1.RPCPort+v1.DefaultPortMargin)},
				"p2p":         map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", v1.P2PPort+v1.DefaultPortMargin)},
				"pprof_laddr": fmt.Sprintf("0.0.0.0:%d", v1.PPROFPort+v1.DefaultPortMargin),
			},
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
		ExpectedFirstValidator: v1.Validator{
			Name:   "user1",
			Bonded: "100000000stake",
			App: map[string]interface{}{
				"grpc":     map[interface{}]interface{}{"address": "localhost:8080"},
				"grpc-web": map[interface{}]interface{}{"address": "localhost:80802"},
				"api":      map[interface{}]interface{}{"address": "localhost:80801"},
			},
			Config: map[string]interface{}{
				"rpc":         map[interface{}]interface{}{"laddr": "localhost:80807"},
				"p2p":         map[interface{}]interface{}{"laddr": "localhost:80804"},
				"pprof_laddr": "localhost:80809",
			},
		},
		ExpectedSecondValidator: v1.Validator{
			Name:   "user2",
			Bonded: "100000000stake",
			App: map[string]interface{}{
				"grpc":     map[interface{}]interface{}{"address": "localhost:8180"},
				"grpc-web": map[interface{}]interface{}{"address": "localhost:81802"},
				"api":      map[interface{}]interface{}{"address": "localhost:81801"},
			},
			Config: map[string]interface{}{
				"rpc":         map[interface{}]interface{}{"laddr": "localhost:81807"},
				"p2p":         map[interface{}]interface{}{"laddr": "localhost:81804"},
				"pprof_laddr": "localhost:81809",
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			conf, err := chainconfig.Parse(strings.NewReader(test.Input))
			require.NoError(t, err)
			require.Equal(t, config.Version(1), conf.GetVersion())
			require.Equal(t, test.ExpectedFirstValidator, conf.Validators[0])
			require.Equal(t, test.ExpectedSecondValidator, conf.Validators[1])
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
validators:
  - name: user1
    bonded: "100000000stake"
    app:
      grpc:
        address: "localhost:8080"
      api:
        address: "localhost:80801"
faucet:
  host: "0.0.0.0:4600"
  port: 4700
`

	conf, err := chainconfig.Parse(strings.NewReader(confyml))
	assert.Nil(t, err)
	require.Equal(t, 1, len(conf.Validators))
	validator := conf.Validators[0]
	require.Equal(t, "localhost:8080", validator.GetGRPC())
	require.Equal(t, "localhost:80801", validator.GetAPI())
	require.Equal(t, "100000000stake", validator.Bonded)
	require.Equal(t, "user1", validator.Name)
}

func TestIsConfigLatest(t *testing.T) {
	path := "testdata/configv0.yaml"
	version, latest, err := chainconfig.IsConfigLatest(path)
	require.NoError(t, err)
	require.Equal(t, false, latest)
	require.Equal(t, config.Version(0), version)

	path = "testdata/configv1.yaml"
	version, latest, err = chainconfig.IsConfigLatest(path)
	require.NoError(t, err)
	require.Equal(t, true, latest)
	require.Equal(t, chainconfig.LatestVersion, version)
}

func TestMigrateLatest(t *testing.T) {
	sourceFile := "testdata/configv0.yaml"
	tempFile := "testdata/temp.yaml"
	input, err := ioutil.ReadFile(sourceFile)
	require.NoError(t, err)

	err = ioutil.WriteFile(tempFile, input, 0644)
	require.NoError(t, err)

	err = MigrateLatest(tempFile)
	require.NoError(t, err)

	targetFile, err := chainconfig.ParseFile(tempFile)
	require.NoError(t, err)

	expectedFile, err := chainconfig.ParseFile("testdata/configv1.yaml")
	require.NoError(t, err)

	require.Equal(t, expectedFile, targetFile)

	// Remove the temp file
	require.NoError(t, os.Remove(tempFile))
}

func TestValidatorWithGentx(t *testing.T) {
	confyml := `
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
        address: "localhost:8080"
      api:
        address: "localhost:80801"
    gentx:
      amount: 1000stake
      moniker: mymoniker
      commission-max-change-rate: max-rate
      home: /path/to/home/dir
      keyring-backend: os
      chain-id: test-chain-1
      commission-max-rate: 1.0
      commission-rate: 0.07
      details: no details
      security-contact: no security-contact
      website: no website
      gas-prices: GasPrices-1
      identity: Identity-1
      min-self-delegation: MinSelfDelegation-1
`

	conf, err := chainconfig.Parse(strings.NewReader(confyml))
	assert.Nil(t, err)
	require.Equal(t, 1, len(conf.Validators))
	validator := conf.Validators[0]
	gentx := validator.Gentx
	require.Equal(t, "1000stake", gentx.Amount)
	require.Equal(t, "mymoniker", gentx.Moniker)
	require.Equal(t, "max-rate", gentx.CommissionMaxChangeRate)
	require.Equal(t, "/path/to/home/dir", gentx.Home)
	require.Equal(t, "os", gentx.KeyringBackend)
	require.Equal(t, "test-chain-1", gentx.ChainID)
	require.Equal(t, "1.0", gentx.CommissionMaxRate)
	require.Equal(t, "0.07", gentx.CommissionRate)
	require.Equal(t, "no details", gentx.Details)
	require.Equal(t, "no security-contact", gentx.SecurityContact)
	require.Equal(t, "no website", gentx.Website)
	require.Equal(t, "GasPrices-1", gentx.GasPrices)
	require.Equal(t, "MinSelfDelegation-1", gentx.MinSelfDelegation)
	require.Equal(t, "Identity-1", gentx.Identity)
}
