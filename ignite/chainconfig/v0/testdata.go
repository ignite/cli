package v0

import (
	"github.com/ignite-hq/cli/ignite/chainconfig/common"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

var v = "bob"

func GetInitialV0Config() common.Config {
	return &Config{
		Validator: Validator{
			Name:   "alice",
			Staked: "100000000stake",
		},
		Init: common.Init{
			App:    nil,
			Client: nil,
			Config: nil,
		},
		Host: common.Host{
			// when in Docker on MacOS, it only works with 0.0.0.0.
			RPC:     "localhost:53803",
			P2P:     "localhost:50198",
			Prof:    "localhost:53030",
			GRPC:    "localhost:53831",
			GRPCWeb: "localhost:46531",
			API:     "localhost:51028",
		},
		BaseConfig: common.BaseConfig{
			Version: 0,
			Build: common.Build{
				Proto: common.Proto{
					Path: "proto",
					ThirdPartyPaths: []string{
						"third_party/proto",
						"proto_vendor",
					},
				},
			},
			Faucet: common.Faucet{
				Name:     &v,
				Host:     "0.0.0.0:4500",
				Coins:    []string{"20000token", "200000000stake"},
				CoinsMax: []string{"10000token", "100000000stake"},
				Port:     48772,
			},
			Accounts: []common.Account{
				{
					Name:  "alice",
					Coins: []string{"20000token", "200000000stake"},
				},
				{
					Name:  "bob",
					Coins: []string{"10000token", "100000000stake"},
				},
			},
			Client: common.Client{
				Vuex: common.Vuex{
					Path: "vue/src/store",
				},
				Dart: common.Dart{
					Path: "",
				},
				OpenAPI: common.OpenAPI{
					Path: "docs/static/openapi.yml",
				},
			},
			Genesis: map[string]interface{}{},
		},
	}
}

func GetConvertedLatestConfig() common.Config {

	return &v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
				Client: nil,
				App: map[string]interface{}{"grpc": map[string]interface{}{"address": "localhost:53831"},
					"grpc-web": map[string]interface{}{"address": "localhost:46531"}, "api": map[string]interface{}{"address": "localhost:51028"}},
				Config: map[string]interface{}{"rpc": map[string]interface{}{"laddr": "localhost:53803"},
					"p2p": map[string]interface{}{"laddr": "localhost:50198"}, "pprof_laddr": "localhost:53030"},
			},
		},
		BaseConfig: common.BaseConfig{
			Version: 1,
			Build: common.Build{
				Proto: common.Proto{
					Path: "proto",
					ThirdPartyPaths: []string{
						"third_party/proto",
						"proto_vendor",
					},
				},
			},
			Faucet: common.Faucet{
				Name:     &v,
				Host:     "0.0.0.0:4500",
				Coins:    []string{"20000token", "200000000stake"},
				CoinsMax: []string{"10000token", "100000000stake"},
				Port:     48772,
			},
			Accounts: []common.Account{
				{
					Name:  "alice",
					Coins: []string{"20000token", "200000000stake"},
				},
				{
					Name:  "bob",
					Coins: []string{"10000token", "100000000stake"},
				},
			},
			Client: common.Client{
				Vuex: common.Vuex{
					Path: "vue/src/store",
				},
				Dart: common.Dart{
					Path: "",
				},
				OpenAPI: common.OpenAPI{
					Path: "docs/static/openapi.yml",
				},
			},
			Genesis: map[string]interface{}{},
		},
	}
}
