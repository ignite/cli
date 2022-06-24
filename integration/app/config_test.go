//go:build !relayer
// +build !relayer

package app_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/confile"
	"github.com/ignite/cli/ignite/pkg/randstr"
	envtest "github.com/ignite/cli/integration"
)

func TestOverwriteSDKConfigsAndChainID(t *testing.T) {
	var (
		env               = envtest.New(t)
		appname           = randstr.Runes(10)
		path              = env.Scaffold(fmt.Sprintf("github.com/test/%s", appname))
		servers           = env.RandomizeServerPorts(path, "")
		ctx, cancel       = context.WithCancel(env.Ctx())
		isBackendAliveErr error
	)

	var c chainconfig.Config

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
	env.Must(env.Serve("should serve", path, "", "", envtest.ExecCtx(ctx)))
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
		cf := confile.New(c.ec, filepath.Join(env.AppdHome(appname), c.relpath))
		require.NoError(t, cf.Load(&conf))
		require.Equal(t, c.expectedVal, conf[c.key])
	}
}
