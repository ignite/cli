package plugininternal

import (
	"context"
	"strconv"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

// TODO move me to ignite org
const consumerPlugin = "github.com/tbruyelle/cli-plugin-consumer"

// ConsumerWriteGenesis writes validators in the consumer module genesis.
// NOTE(tb): Using a plugin for this task avoids having the interchain-security
// dependency in Ignite.
func ConsumerWriteGenesis(ctx context.Context, c plugin.Chainer) error {
	_, err := Execute(ctx, consumerPlugin, []string{"writeGenesis"}, plugin.WithChain(c))
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
	out, err := Execute(ctx, consumerPlugin, []string{"isInitialized"}, plugin.WithChain(c))
	if err != nil {
		return false, errors.Errorf("execute consumer plugin 'isInitialized': %w", err)
	}
	b, err := strconv.ParseBool(out)
	if err != nil {
		return false, errors.Errorf("invalid consumer plugin 'isInitialized' output, got '%s': %w", out, err)
	}
	return b, nil
}
