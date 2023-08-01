package events_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/events"
)

func TestNew(t *testing.T) {
	msg := "message"
	cases := []struct {
		name, message       string
		inProgress, hasIcon bool
		options             []events.Option
		event               events.Event
	}{
		{
			name:  "event",
			event: events.Event{},
		},
		{
			name:       "event start",
			message:    msg,
			inProgress: true,
			options:    []events.Option{events.ProgressStart()},
			event:      events.New(msg, events.ProgressStart()),
		},
		{
			name:       "event update",
			message:    msg,
			inProgress: true,
			options:    []events.Option{events.ProgressUpdate()},
			event:      events.New(msg, events.ProgressUpdate()),
		},
		{
			name:    "event finish",
			message: msg,
			options: []events.Option{events.ProgressFinish()},
			event:   events.New(msg, events.ProgressFinish()),
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			e := events.New(tt.message, tt.options...)

			// Assert
			require.Equal(t, tt.event, e)
			require.Equal(t, tt.inProgress, e.InProgress())
		})
	}
}
