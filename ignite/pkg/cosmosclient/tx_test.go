package cosmosclient_test

import (
	"encoding/json"
	"testing"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/stretchr/testify/require"
)

func TestUnmarshallEvents(t *testing.T) {
	// Arrange
	wantEvents := []cosmosclient.TXEvent{
		{
			Type: "test",
			Attributes: []cosmosclient.TXEventAttribute{
				{Key: "foo", Value: "bar"},
				{Key: "baz", Value: 42.0},
			},
		},
	}

	log := []struct {
		Events []cosmosclient.TXEvent `json:"events"`
	}{
		{Events: wantEvents},
	}

	raw, err := json.Marshal(log)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	events, err := cosmosclient.UnmarshallEvents(raw)

	// Assert
	require.NoError(t, err)
	require.Equal(t, wantEvents, events)
}
