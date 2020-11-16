// +build !relayer

package integration_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/randstr"
	"github.com/tendermint/starport/starport/services/chain/conf"
)

func TestOverwriteSDKConfigsAndChainID(t *testing.T) {
	t.Parallel()

	for _, sdkVersion := range []string{Launchpad, Stargate} {
		sdkVersion := sdkVersion
		t.Run(sdkVersion, func(t *testing.T) {
			t.Parallel()

			testOverwriteSDKConfigsAndChainID(t, sdkVersion)
		})
	}
}

func testOverwriteSDKConfigsAndChainID(t *testing.T, sdkVersion string) {
	var (
		env               = newEnv(t)
		appname           = randstr.Runes(10)
		path              = env.Scaffold(appname, sdkVersion)
		servers           = env.RandomizeServerPorts(path)
		ctx, cancel       = context.WithCancel(env.Ctx())
		isBackendAliveErr error
	)

	var c conf.Config

	cf := confile.New(confile.DefaultYAMLEncodingCreator, filepath.Join(path, "config.yml"))
	require.NoError(t, cf.Load(&c))

	c.Genesis = map[string]interface{}{"chain_id": "cosmos"}
	c.Init.App = map[string]interface{}{"hello": "cosmos"}
	c.Init.Config = map[string]interface{}{"fast_sync": false}

	require.NoError(t, cf.Save(c))

	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve", path, ExecCtx(ctx)))
	require.NoError(t, isBackendAliveErr, "app cannot get online in time")

	configs := []struct {
		ec          confile.EncodingCreator
		relpath     string
		key         string
		expectedVal interface{}
	}{
		{confile.DefaultJSONEncodingCreator, "config/genesis.json", "chain_id", "cosmos"},
		{confile.DefaultTOMLEncodingCreator, "config/app.toml", "hello", "cosmos"},
		{confile.DefaultTOMLEncodingCreator, "config/config.toml", "fast_sync", false},
	}

	for _, c := range configs {
		var conf map[string]interface{}
		cf := confile.New(c.ec, filepath.Join(env.AppdHome(appname, sdkVersion), c.relpath))
		require.NoError(t, cf.Load(&conf))
		require.Equal(t, c.expectedVal, conf[c.key])
	}
}
