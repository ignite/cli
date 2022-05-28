package v1

import (
	"testing"

	"github.com/ignite-hq/cli/ignite/chainconfig/common"
	"github.com/stretchr/testify/require"
)

func TestListValidators(t *testing.T) {
	conf := Config{
		Validators: []Validator{
			{
				Name:   "test-name",
				Bonded: "101ATOM",
			},
			{
				Name:   "test-name-1",
				Bonded: "102ATOM",
			},
		},
	}

	require.Equal(t, []common.Validator{
		&Validator{
			Name:   "test-name",
			Bonded: "101ATOM",
		}, &Validator{
			Name:   "test-name-1",
			Bonded: "102ATOM",
		}}, conf.ListValidators())
}

func TestGetHost(t *testing.T) {
	conf := Config{
		Validators: []Validator{
			{
				Name:   "test-name-1",
				Bonded: "102ATOM",
				App: map[string]interface{}{"grpc": map[string]interface{}{"address": "0.0.0.0:9090"},
					"grpc-web": map[string]interface{}{"address": "0.0.0.0:9091"}, "api": map[string]interface{}{"address": "0.0.0.0:1317"}},
				Config: map[string]interface{}{"rpc": map[string]interface{}{"laddr": "0.0.0.0:26657"},
					"p2p": map[string]interface{}{"laddr": "0.0.0.0:26656"}, "pprof_laddr": "0.0.0.0:6060"},
			},
		},
	}

	require.Equal(t, common.Host{
		RPC:     "0.0.0.0:26657",
		P2P:     "0.0.0.0:26656",
		Prof:    "0.0.0.0:6060",
		GRPC:    "0.0.0.0:9090",
		GRPCWeb: "0.0.0.0:9091",
		API:     "0.0.0.0:1317",
	}, conf.GetHost())
}
