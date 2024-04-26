package plugininternal

import (
	"context"
	"strconv"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

const (
	PluginConsumerVersion = "36532a375ecec4a228dc9b320c1576a2a1fafd34"
	PluginConsumerPath    = "github.com/ignite/apps/consumer@" + PluginConsumerVersion
)

// ConsumerWriteGenesis writes validators in the consumer module genesis.
// NOTE(tb): Using a plugin for this task avoids having the interchain-security
// dependency in Ignite.
func ConsumerWriteGenesis(ctx context.Context, c plugin.Chainer) error {
	_, err := Execute(ctx, PluginConsumerPath, []string{"writeGenesis"}, plugin.WithChain(c))
	if err != nil {
		return errors.Errorf("execute consumer plugin 'writeGenesis': %w", err)
	}
	return nil
}

// ConsumerIsInitialized returns true if the consumer chain's genesis c has
// a consumer module entry with an initial validator set.
// NOTE(tb): Using a plugin for this task avoids having the interchain-security
// dependency in Ignite.
func ConsumerIsInitialized(ctx context.Context, c plugin.Chainer) (bool, error) {
	out, err := Execute(ctx, PluginConsumerPath, []string{"isInitialized"}, plugin.WithChain(c))
	if err != nil {
		return false, errors.Errorf("execute consumer plugin 'isInitialized': %w", err)
	}
	b, err := strconv.ParseBool(out)
	if err != nil {
		return false, errors.Errorf("invalid consumer plugin 'isInitialized' output, got '%s': %w", out, err)
	}
	return b, nil
}
