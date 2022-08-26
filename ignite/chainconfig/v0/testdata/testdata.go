package testdata

import (
	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v0 "github.com/ignite-hq/cli/ignite/chainconfig/v0"
)

// GetConfigV0 returns a v0 config initialized with random values.
func GetConfigV0() *v0.Config {
	faucetName := "bob"

	return &v0.Config{
		BaseConfig: config.BaseConfig{
			Version: 0,
			Build: config.Build{
				Proto: config.Proto{
					Path: "proto",
					ThirdPartyPaths: []string{
						"third_party/proto",
						"proto_vendor",
					},
				},
			},
			Faucet: config.Faucet{
				Name:     &faucetName,
				Host:     "0.0.0.0:4500",
				Coins:    []string{"20000token", "200000000stake"},
				CoinsMax: []string{"10000token", "100000000stake"},
				Port:     48772,
			},
			Accounts: []config.Account{
				{
					Name:  "alice",
					Coins: []string{"20000token", "200000000stake"},
				},
				{
					Name:  "bob",
					Coins: []string{"10000token", "100000000stake"},
				},
			},
			Client: config.Client{
				Vuex: config.Vuex{
					Path: "vue/src/store",
				},
				Dart: config.Dart{
					Path: "",
				},
				OpenAPI: config.OpenAPI{
					Path: "docs/static/openapi.yml",
				},
			},
			Genesis: map[string]interface{}{},
		},
		Validator: v0.Validator{
			Name:   "alice",
			Staked: "100000000stake",
		},
		Init: config.Init{
			App:    map[string]interface{}{"test-app": "test-app"},
			Config: map[string]interface{}{"test-config": "test-config"},
			Client: map[string]interface{}{"test-client": "test-client"},
		},
		Host: config.Host{
			RPC:     "localhost:53803",
			P2P:     "localhost:50198",
			Prof:    "localhost:53030",
			GRPC:    "localhost:53831",
			GRPCWeb: "localhost:46531",
			API:     "localhost:51028",
		},
	}
}
