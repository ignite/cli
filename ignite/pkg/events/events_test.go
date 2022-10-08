package events_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/events"
)

func TestEventIsOngoing(t *testing.T) {
	tests := []struct {
		name    string
		status  events.ProgressIndication
		message string
		want    bool
	}{
		{"status ongoing", events.IndicationStart, "description", true},
		{"status done", events.IndicationFinish, "description", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := events.Event{
				ProgressIndication: tt.status,
				Message:            tt.message,
			}

			require.Equal(t, tt.want, e.IsOngoing())
		})
	}
}

func TestEventString(t *testing.T) {
	tests := []struct {
		name    string
		status  events.ProgressIndication
		message string
		want    string
	}{
		{
			name:    "status done",
			status:  events.IndicationFinish,
			message: "message",
			want:    "message",
		},
		{
			name:    "status ongoing",
			status:  events.IndicationStart,
			message: "message",
			want:    "message",
		},
		{
			name:    "status ongoing with empty message",
			status:  events.IndicationStart,
			message: "",
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := events.Event{
				ProgressIndication: tt.status,
				Message:            tt.message,
			}

			require.Equal(t, tt.want, e.Message)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		message string
		options []events.Option
		want    events.Event
	}{
		{
			name: "zero value args",
			want: events.Event{},
		},
		{
			name:    "status ongoing",
			options: []events.Option{events.ProgressStarted()},
			message: "message",
			want:    events.Event{ProgressIndication: 1, Message: "message"},
		},
		{
			name:    "status done",
			options: []events.Option{events.ProgressFinished()},
			message: "message",
			want:    events.Event{ProgressIndication: 2, Message: "message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, events.New(tt.message, tt.options...))
		})
	}
}
