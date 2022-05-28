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

	require.Equal(t, []*Validator{
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

func TestGetAddressPort(t *testing.T) {
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

	require.Equal(t, 1, len(conf.Validators))
	validator := conf.Validators[0]
	require.Equal(t, "0.0.0.0", validator.GetGRPCAddress())
	require.Equal(t, 9090, validator.GetGRPCPort())
	require.Equal(t, "0.0.0.0", validator.GetP2PAddress())
	require.Equal(t, 26656, validator.GetP2PPort())
	require.Equal(t, "0.0.0.0", validator.GetGRPCWebAddress())
	require.Equal(t, 9091, validator.GetGRPCWebPort())
	require.Equal(t, "0.0.0.0", validator.GetAPIAddress())
	require.Equal(t, 1317, validator.GetAPIPort())
	require.Equal(t, "0.0.0.0", validator.GetRPCAddress())
	require.Equal(t, 26657, validator.GetRPCPort())
	require.Equal(t, "0.0.0.0", validator.GetProfAddress())
	require.Equal(t, 6060, validator.GetProfPort())

	validator = validator.IncreasePort(10)
	require.Equal(t, "0.0.0.0", validator.GetGRPCAddress())
	require.Equal(t, 9100, validator.GetGRPCPort())
	require.Equal(t, "0.0.0.0", validator.GetP2PAddress())
	require.Equal(t, 26666, validator.GetP2PPort())
	require.Equal(t, "0.0.0.0", validator.GetGRPCWebAddress())
	require.Equal(t, 9101, validator.GetGRPCWebPort())
	require.Equal(t, "0.0.0.0", validator.GetAPIAddress())
	require.Equal(t, 1327, validator.GetAPIPort())
	require.Equal(t, "0.0.0.0", validator.GetRPCAddress())
	require.Equal(t, 26667, validator.GetRPCPort())
	require.Equal(t, "0.0.0.0", validator.GetProfAddress())
	require.Equal(t, 6070, validator.GetProfPort())
}
