package cosmosclient

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshallEvents(t *testing.T) {
	// Arrange
	wantEvents := []TXEvent{
		{
			Type: "test",
			Attributes: []TXEventAttribute{
				{Key: "foo", Value: "bar"},
				{Key: "baz", Value: 42.0},
			},
		},
	}

	log := []struct {
		Events []TXEvent `json:"events"`
	}{
		{Events: wantEvents},
	}

	raw, err := json.Marshal(log)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	events, err := UnmarshallEvents(raw)

	// Assert
	require.NoError(t, err)
	require.Equal(t, wantEvents, events)
}
