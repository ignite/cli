package v0_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v0testdata "github.com/ignite-hq/cli/ignite/chainconfig/v0/testdata"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

func TestConvertV0ToV1(t *testing.T) {
	// Arrange
	cfgV0 := v0testdata.GetConfig(t)

	// Act
	c, err := cfgV0.ConvertNext()
	cfgV1, ok := c.(*v1.Config)

	// Assert
	require.NoError(t, err)
	require.Truef(t, ok, "expected *v1.Config, got %T", c)
	require.Equal(t, config.Version(1), cfgV1.GetVersion())

	// Assert: BaseConfig values
	require.Equal(t, cfgV0.Build, cfgV1.Build)
	require.Equal(t, cfgV0.Accounts, cfgV1.Accounts)
	require.Equal(t, cfgV0.Faucet, cfgV1.Faucet)
	require.Equal(t, cfgV0.Client, cfgV1.Client)
	require.Equal(t, cfgV0.Genesis, cfgV1.Genesis)

	// Assert: Values that must be migrated to the first validator
	require.Len(t, cfgV1.Validators, 1)
	require.Equal(t, cfgV0.Validator.Name, cfgV1.Validators[0].Name)
	require.Equal(t, cfgV0.Validator.Staked, cfgV1.Validators[0].Bonded)
	require.Equal(t, cfgV0.Init.Home, cfgV1.Validators[0].Home)
	require.Equal(t, cfgV0.Init.KeyringBackend, cfgV1.Validators[0].KeyringBackend)
	require.Equal(t, cfgV0.Init.Client, cfgV1.Validators[0].Client)
	require.Equal(t, cfgV0.Host.RPC, cfgV1.Validators[0].GetRPC())
	require.Equal(t, cfgV0.Host.P2P, cfgV1.Validators[0].GetP2P())
	require.Equal(t, cfgV0.Host.GRPC, cfgV1.Validators[0].GetGRPC())
	require.Equal(t, cfgV0.Host.GRPCWeb, cfgV1.Validators[0].GetGRPCWeb())
	require.Equal(t, cfgV0.Host.Prof, cfgV1.Validators[0].GetProf())
	require.Equal(t, cfgV0.Host.API, cfgV1.Validators[0].GetAPI())
}
