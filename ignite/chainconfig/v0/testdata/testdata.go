package testdata

import (
	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v0 "github.com/ignite-hq/cli/ignite/chainconfig/v0"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

var faucetName = "bob"

func GetInitialV0Config() *v0.Config {
	return &v0.Config{
		BaseConfig: config.BaseConfig{
			ConfigVersion: 0,
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
			// when in Docker on MacOS, it only works with 0.0.0.0.
			RPC:     "localhost:53803",
			P2P:     "localhost:50198",
			Prof:    "localhost:53030",
			GRPC:    "localhost:53831",
			GRPCWeb: "localhost:46531",
			API:     "localhost:51028",
		},
	}
}

func GetConvertedLatestConfig() *v1.Config {
	return &v1.Config{
		BaseConfig: config.BaseConfig{
			ConfigVersion: 1,
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
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
				Client: map[string]interface{}{"test-client": "test-client"},
				App: map[string]interface{}{
					"grpc":     map[string]interface{}{"address": "localhost:53831"},
					"grpc-web": map[string]interface{}{"address": "localhost:46531"},
					"api":      map[string]interface{}{"address": "localhost:51028"},
					"test-app": "test-app",
				},
				Config: map[string]interface{}{
					"rpc":         map[string]interface{}{"laddr": "localhost:53803"},
					"p2p":         map[string]interface{}{"laddr": "localhost:50198"},
					"pprof_laddr": "localhost:53030",
					"test-config": "test-config",
				},
			},
		},
	}
}
