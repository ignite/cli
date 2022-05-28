//go:build !relayer
// +build !relayer

package app_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/imdario/mergo"

	"github.com/stretchr/testify/require"

	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
	"github.com/ignite-hq/cli/ignite/pkg/confile"
	"github.com/ignite-hq/cli/ignite/pkg/randstr"
	envtest "github.com/ignite-hq/cli/integration"
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

	var c v1.Config

	cf := confile.New(confile.DefaultYAMLEncodingCreator, filepath.Join(path, "config.yml"))
	require.NoError(t, cf.Load(&c))

	c.Genesis = map[string]interface{}{"chain_id": "cosmos"}
	defaultValidator := v1.Validator{
		App:    map[string]interface{}{"hello": "cosmos"},
		Config: map[string]interface{}{"fast_sync": false},
	}

	for i := range c.Validators {
		err := mergo.Merge(&c.Validators[i], defaultValidator)
		require.NoError(t, err)
	}

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
