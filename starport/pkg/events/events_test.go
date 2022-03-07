package events

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBusSend(t *testing.T) {
	tests := []struct {
		name  string
		bus   Bus
		event Event
	}{
		{
			name: "send status ongoing event",
			bus:  make(Bus),
			event: Event{
				status:      StatusOngoing,
				Description: "description",
			},
		},
		{
			name: "send status done event",
			bus:  make(Bus),
			event: Event{
				status:      StatusDone,
				Description: "description",
			},
		},
		{
			name: "send event on nil bus",
			bus:  nil,
			event: Event{
				status:      StatusDone,
				Description: "description",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.bus.Send(tt.event)
			if tt.bus != nil {
				require.Equal(t, tt.event, <-tt.bus)
			}
			tt.bus.Shutdown()
		})
	}
}

func TestBusShutdown(t *testing.T) {
	tests := []struct {
		name string
		bus  Bus
	}{
		{
			name: "shutdown nil bus",
			bus:  nil,
		},
		{
			name: "shutdown bus correctly",
			bus:  make(Bus),
		},
		{
			name: "shutdown bus with size correctly",
			bus:  make(Bus, 1),
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
		status      Status
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"status ongoing", fields{StatusOngoing, "description"}, true},
		{"status done", fields{StatusDone, "description"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Event{
				status:      tt.fields.status,
				Description: tt.fields.Description,
			}
			require.Equal(t, tt.want, e.IsOngoing())
		})
	}
}

func TestEventText(t *testing.T) {
	type fields struct {
		status      Status
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "status done",
			fields: fields{StatusDone, "description"},
			want:   "description",
		},
		{
			name:   "status ongoing",
			fields: fields{StatusOngoing, "description"},
			want:   "description...",
		},
		{
			name:   "status ongoing with empty description",
			fields: fields{StatusOngoing, ""},
			want:   "...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Event{
				status:      tt.fields.status,
				Description: tt.fields.Description,
			}
			require.Equal(t, tt.want, e.Text())
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		status      Status
		description string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{"zero value args", args{}, Event{}},
		{"large value args", args{99999, "description"}, Event{99999, "description"}},
		{"status ongoing", args{StatusOngoing, "description"}, Event{0, "description"}},
		{"status done", args{StatusDone, "description"}, Event{1, "description"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, New(tt.args.status, tt.args.description))
		})
	}
}

func TestNewBus(t *testing.T) {
	tests := []struct {
		name  string
		event Event
	}{
		{"new bus with status done event", Event{StatusDone, "description"}},
		{"new bus with status ongoing event", Event{StatusOngoing, "description"}},
		{"new bus with zero value event", Event{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewBus()
			defer bus.Shutdown()
			for i := 0; i < 10; i++ {
				go bus.Send(tt.event)
				require.Equal(t, tt.event, <-bus)
			}
		})
	}
}
