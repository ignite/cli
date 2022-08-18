package v1_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

func TestValidators(t *testing.T) {
	conf := v1.Config{
		Validators: []v1.Validator{
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

	require.Equal(t, []v1.Validator{
		{
			Name:   "test-name",
			Bonded: "101ATOM",
		}, {
			Name:   "test-name-1",
			Bonded: "102ATOM",
		},
	}, conf.Validators)
}

func TestGetAddressPort(t *testing.T) {
	conf := v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "test-name-1",
				Bonded: "102ATOM",
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

func TestClone(t *testing.T) {
	config := &v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
			},
		},
	}
	clone := config.Clone()
	require.Equal(t, config, clone)

	clone.(*v1.Config).Validators = []v1.Validator{
		{
			Name:   "test",
			Bonded: "stakedvalue",
		},
	}
	require.NotEqual(t, config, clone)
	require.Equal(t, []v1.Validator{
		{
			Name:   "test",
			Bonded: "stakedvalue",
		},
	}, clone.(*v1.Config).Validators)
}

func TestChangeValidators(t *testing.T) {
	conf := v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "test-name-1",
				Bonded: "102ATOM",
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
		},
	}

	require.Equal(t, "", conf.Validators[0].Home)
	require.Equal(t, "test-name-1", conf.Validators[0].Name)

	conf.Validators[0].Home = "test-home"
	conf.Validators[0].Name = "test-name"
	require.Equal(t, "test-home", conf.Validators[0].Home)
	require.Equal(t, "test-name", conf.Validators[0].Name)
}

func TestFillValidatorsDefaults(t *testing.T) {
	tests := []struct {
		TestName         string
		InputConf        v1.Config
		DefaultValidator v1.Validator
		ExpectedConf     v1.Config
	}{{
		TestName: "Config contains the validator with the ports defined",
		InputConf: v1.Config{
			Validators: []v1.Validator{
				{
					Name:   "test-name-1",
					Bonded: "102ATOM",
					App: map[string]interface{}{
						"grpc":     map[string]interface{}{"address": "0.0.0.0:19090"},
						"grpc-web": map[string]interface{}{"address": "0.0.0.0:19091"},
						"api":      map[string]interface{}{"address": "0.0.0.0:2317"},
					},
					Config: map[string]interface{}{
						"rpc":         map[string]interface{}{"laddr": "0.0.0.0:36657"},
						"p2p":         map[string]interface{}{"laddr": "0.0.0.0:36656"},
						"pprof_laddr": "0.0.0.0:7060",
					},
				},
				{
					Name:   "test-name-2",
					Bonded: "103ATOM",
				},
				{
					Name:   "test-name-3",
					Bonded: "104ATOM",
				},
			},
		},
		DefaultValidator: v1.Validator{
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
		ExpectedConf: v1.Config{
			Validators: []v1.Validator{
				{
					Name:   "test-name-1",
					Bonded: "102ATOM",
					App: map[string]interface{}{
						"grpc":     map[string]interface{}{"address": "0.0.0.0:19090"},
						"grpc-web": map[string]interface{}{"address": "0.0.0.0:19091"},
						"api":      map[string]interface{}{"address": "0.0.0.0:2317"},
					},
					Config: map[string]interface{}{
						"rpc":         map[string]interface{}{"laddr": "0.0.0.0:36657"},
						"p2p":         map[string]interface{}{"laddr": "0.0.0.0:36656"},
						"pprof_laddr": "0.0.0.0:7060",
					},
				},
				{
					Name:   "test-name-2",
					Bonded: "103ATOM",
					App: map[string]interface{}{
						"grpc":     map[string]interface{}{"address": "0.0.0.0:19100"},
						"grpc-web": map[string]interface{}{"address": "0.0.0.0:19101"},
						"api":      map[string]interface{}{"address": "0.0.0.0:2327"},
					},
					Config: map[string]interface{}{
						"rpc":         map[string]interface{}{"laddr": "0.0.0.0:36667"},
						"p2p":         map[string]interface{}{"laddr": "0.0.0.0:36666"},
						"pprof_laddr": "0.0.0.0:7070",
					},
				},
				{
					Name:   "test-name-3",
					Bonded: "104ATOM",
					App: map[string]interface{}{
						"grpc":     map[string]interface{}{"address": "0.0.0.0:19110"},
						"grpc-web": map[string]interface{}{"address": "0.0.0.0:19111"},
						"api":      map[string]interface{}{"address": "0.0.0.0:2337"},
					},
					Config: map[string]interface{}{
						"rpc":         map[string]interface{}{"laddr": "0.0.0.0:36677"},
						"p2p":         map[string]interface{}{"laddr": "0.0.0.0:36676"},
						"pprof_laddr": "0.0.0.0:7080",
					},
				},
			},
		},
	}, {
		TestName: "Config contains the validator with the ports undefined",
		InputConf: v1.Config{
			Validators: []v1.Validator{
				{
					Name:   "test-name-1",
					Bonded: "102ATOM",
				},
				{
					Name:   "test-name-2",
					Bonded: "103ATOM",
				},
			},
		},
		DefaultValidator: v1.Validator{
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
		ExpectedConf: v1.Config{
			Validators: []v1.Validator{
				{
					Name:   "test-name-1",
					Bonded: "102ATOM",
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
				{
					Name:   "test-name-2",
					Bonded: "103ATOM",
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
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			test.InputConf.FillValidatorsDefaults(test.DefaultValidator)
			require.Equal(t, test.ExpectedConf, test.InputConf)
		})
	}
}
