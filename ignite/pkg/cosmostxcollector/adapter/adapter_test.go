package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/v29/ignite/pkg/cosmostxcollector/query"
)

type testAdapter struct{}

func (testAdapter) Save(context.Context, []cosmosclient.TX) error  { return nil }
func (testAdapter) GetType() string                                { return "test" }
func (testAdapter) Init(context.Context) error                     { return nil }
func (testAdapter) GetLatestHeight(context.Context) (int64, error) { return 1, nil }
func (testAdapter) QueryEvents(context.Context, query.EventQuery) ([]query.Event, error) {
	return nil, nil
}

func (testAdapter) Query(context.Context, query.Query) (query.Cursor, error) {
	return nil, nil
}

func TestTestAdapterImplementsInterfaces(t *testing.T) {
	var _ Saver = testAdapter{}
	var a Adapter = testAdapter{}

	require.Equal(t, "test", a.GetType())
	require.NoError(t, a.Save(context.Background(), nil))
}
