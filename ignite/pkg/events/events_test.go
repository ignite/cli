package events_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/pkg/events"
)

func TestBusSend(t *testing.T) {
	tests := []struct {
		name    string
		bus     events.Bus
		event   events.Event
		options []events.Option
	}{
		{
			name:  "send status ongoing event",
			bus:   events.NewBus(),
			event: events.New(events.StringContent("description"), events.ProgressStarted()),
			options: []events.Option{
				events.ProgressStarted(),
			},
		},
		{
			name:  "send status done event",
			bus:   events.NewBus(),
			event: events.New(events.StringContent("description"), events.ProgressFinished()),
			options: []events.Option{
				events.ProgressFinished(),
			},
		},
		{
			name:  "send status neutral event",
			bus:   events.NewBus(),
			event: events.New(events.StringContent("description")),
		},
		{
			name:  "send event on nil bus",
			bus:   events.Bus{},
			event: events.New(events.StringContent("description"), events.ProgressFinished()),
			options: []events.Option{
				events.ProgressFinished(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.bus.Send(tt.event.Content, tt.options...)
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
		status      events.ProgressIndication
		description string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"status ongoing", fields{events.IndicationStart, "description"}, true},
		{"status done", fields{events.IndicationFinish, "description"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := events.Event{
				ProgressIndication: tt.fields.status,
				Content:            events.StringContent(tt.fields.description),
			}
			require.Equal(t, tt.want, e.IsOngoing())
		})
	}
}

func TestEventString(t *testing.T) {
	type fields struct {
		status      events.ProgressIndication
		description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "status done",
			fields: fields{
				status:      events.IndicationFinish,
				description: "description",
			},
			want: "description",
		},
		{
			name: "status ongoing",
			fields: fields{
				status:      events.IndicationStart,
				description: "description",
			},
			want: "description",
		},
		{
			name: "status ongoing with empty description",
			fields: fields{
				status:      events.IndicationStart,
				description: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := events.Event{
				ProgressIndication: tt.fields.status,
				Content:            events.StringContent(tt.fields.description),
			}
			require.Equal(t, tt.want, e.Content.String())
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		content events.Content
		options []events.Option
	}
	tests := []struct {
		name string
		args args
		want events.Event
	}{
		{
			"zero value args",
			args{},
			events.Event{},
		},
		{
			"status ongoing",
			args{
				options: []events.Option{events.ProgressStarted()},
				content: events.StringContent("description"),
			},
			events.Event{ProgressIndication: 1, Content: events.StringContent("description")},
		},
		{
			"status done",
			args{
				options: []events.Option{events.ProgressFinished()},
				content: events.StringContent("description"),
			},
			events.Event{ProgressIndication: 2, Content: events.StringContent("description")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, events.New(tt.args.content, tt.args.options...))
		})
	}
}
