package events_test

import (
	"testing"

	"github.com/gookit/color"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/events"
)

func TestBusSend(t *testing.T) {
	tests := []struct {
		name  string
		bus   events.Bus
		event events.Event
	}{
		{
			name: "send status ongoing event",
			bus:  events.NewBus(),
			event: events.Event{
				Status:      events.StatusOngoing,
				Description: "description",
			},
		},
		{
			name: "send status done event",
			bus:  events.NewBus(),
			event: events.Event{
				Status:      events.StatusDone,
				Description: "description",
			},
		},
		{
			name: "send status neutral event",
			bus:  events.NewBus(),
			event: events.Event{
				Status:      events.StatusNeutral,
				Description: "description",
			},
		},
		{
			name: "send event on nil bus",
			bus:  events.Bus{},
			event: events.Event{
				Status: events.StatusDone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.bus.Send(tt.event)
			if tt.bus.Events() != nil {
				require.Equal(t, tt.event, <-tt.bus.Events())
			}
			tt.bus.Shutdown()
		})
	}
}

func TestBusShutdown(t *testing.T) {
	tests := []struct {
		name string
		bus  events.Bus
	}{
		{
			name: "shutdown nil bus",
			bus:  events.Bus{},
		},
		{
			name: "shutdown bus correctly",
			bus:  events.NewBus(),
		},
		{
			name: "shutdown bus with size correctly",
			bus:  events.NewBus(events.WithCustomBufferSize(1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bus.Shutdown()
		})
	}
}

func TestEventIsOngoing(t *testing.T) {
	type fields struct {
		status      events.Status
		description string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"status ongoing", fields{events.StatusOngoing, "description"}, true},
		{"status done", fields{events.StatusDone, "description"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := events.Event{
				Status:      tt.fields.status,
				Description: tt.fields.description,
			}
			require.Equal(t, tt.want, e.IsOngoing())
		})
	}
}

func TestEventText(t *testing.T) {
	type fields struct {
		status      events.Status
		description string
		textColor   color.Color
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "status done",
			fields: fields{
				status:      events.StatusDone,
				description: "description",
				textColor:   color.Red,
			},
			want: "description",
		},
		{
			name: "status ongoing",
			fields: fields{
				status:      events.StatusOngoing,
				description: "description",
				textColor:   color.Red,
			},
			want: "description...",
		},
		{
			name: "status ongoing with empty description",
			fields: fields{
				status:      events.StatusOngoing,
				description: "",
				textColor:   color.Red,
			},
			want: "...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := events.Event{
				Status:      tt.fields.status,
				Description: tt.fields.description,
				TextColor:   tt.fields.textColor,
			}
			require.Equal(t, e.TextColor.Render(tt.want), e.Text())
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		status      events.Status
		description string
	}
	tests := []struct {
		name string
		args args
		want events.Event
	}{
		{"zero value args", args{}, events.Event{}},
		{"large value args", args{status: 99999, description: "description"}, events.Event{Status: 99999, Description: "description"}},
		{"status ongoing", args{status: events.StatusOngoing, description: "description"}, events.Event{Status: 0, Description: "description"}},
		{"status done", args{status: events.StatusDone, description: "description"}, events.Event{Status: 1, Description: "description"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, events.New(tt.args.status, tt.args.description))
		})
	}
}

func TestNewBus(t *testing.T) {
	tests := []struct {
		name  string
		event events.Event
	}{
		{"new bus with status done event", events.Event{Status: events.StatusDone, Description: "description"}},
		{"new bus with status ongoing event", events.Event{Status: events.StatusOngoing, Description: "description"}},
		{"new bus with zero value event", events.Event{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := events.NewBus()
			defer bus.Shutdown()
			for i := 0; i < 10; i++ {
				go bus.Send(tt.event)
				require.Equal(t, tt.event, <-bus.Events())
			}
		})
	}
}
