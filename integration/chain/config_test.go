//go:build !relayer

package chain_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/confile"
	"github.com/ignite/cli/ignite/pkg/randstr"
	envtest "github.com/ignite/cli/integration"
)

func TestOverwriteSDKConfigsAndChainID(t *testing.T) {
	var (
		env               = envtest.New(t)
		appname           = randstr.Runes(10)
		app               = env.Scaffold(fmt.Sprintf("github.com/test/%s", appname))
		servers           = app.RandomizeServerPorts()
		ctx, cancel       = context.WithCancel(env.Ctx())
		isBackendAliveErr error
	)

	var cfg chainconfig.Config
	cf := confile.New(confile.DefaultYAMLEncodingCreator, filepath.Join(app.SourcePath(), "config.yml"))
	require.NoError(t, cf.Load(&cfg))

	cfg.Genesis = map[string]interface{}{"chain_id": "cosmos"}
	cfg.Validators[0].App["hello"] = "cosmos"
	cfg.Validators[0].Config["fast_sync"] = false

	require.NoError(t, cf.Save(cfg))

	go func() {
		defer cancel()

		isBackendAliveErr = env.IsAppServed(ctx, servers.API)
	}()

	env.Must(app.Serve("should serve", envtest.ExecCtx(ctx)))
	require.NoError(t, isBackendAliveErr, "app cannot get online in time")

	cases := []struct {
		ec      confile.EncodingCreator
		relpath string
		key     string
		want    interface{}
	}{
		{confile.DefaultJSONEncodingCreator, "config/genesis.json", "chain_id", "cosmos"},
		{confile.DefaultTOMLEncodingCreator, "config/app.toml", "hello", "cosmos"},
		{confile.DefaultTOMLEncodingCreator, "config/config.toml", "fast_sync", false},
	}

	for _, tt := range cases {
		var conf map[string]interface{}

		path := filepath.Join(env.AppHome(appname), tt.relpath)
		c := confile.New(tt.ec, path)

		require.NoError(t, c.Load(&conf))
		require.Equalf(t, tt.want, conf[tt.key], "unexpected value for %s", tt.relpath)
	}
}
