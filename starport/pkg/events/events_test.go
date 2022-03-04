package events

import (
	"reflect"
	"testing"
)

func TestBus_Send(t *testing.T) {
	type args struct {
		e Event
	}
	tests := []struct {
		name string
		b    Bus
		args args
	}{
		{
			name: "send status ongoing event",
			b:    make(Bus),
			args: args{Event{
				status:      StatusOngoing,
				Description: "description",
			}},
		},
		{
			name: "send status done event",
			b:    make(Bus),
			args: args{Event{
				status:      StatusDone,
				Description: "description",
			}},
		},
		{
			name: "send event on nil bus",
			b:    nil,
			args: args{Event{
				status:      StatusDone,
				Description: "description",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				tt.b.Send(tt.args.e)
			}()
			if tt.b != nil {
				event := <-tt.b
				if !reflect.DeepEqual(event, tt.args.e) {
					t.Errorf("event = %v, want %v", event, tt.args.e)
				}
			}
			tt.b.Shutdown()
		})
	}
}

func TestBus_Shutdown(t *testing.T) {
	tests := []struct {
		name string
		b    Bus
	}{
		{
			name: "shutdown nil bus",
			b:    nil,
		},
		{
			name: "shutdown bus correctly",
			b:    make(Bus),
		},
		{
			name: "shutdown bus with size correctly",
			b:    make(Bus, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Shutdown()
		})
	}
}

func TestEvent_IsOngoing(t *testing.T) {
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
			if got := e.IsOngoing(); got != tt.want {
				t.Errorf("IsOngoing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_Text(t *testing.T) {
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
			if got := e.Text(); got != tt.want {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
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
			if got := New(tt.args.status, tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBus(t *testing.T) {
	tests := []struct {
		name string
		want Bus
	}{
		{"new bus event chan", make(Bus)},
		{"new nus event chan matches event chan", make(chan Event)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBus()
			if reflect.TypeOf(got).Kind() != reflect.Chan {
				t.Errorf("NewBus() = %v is not a channel", got)
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewBus() = %v, want %v", got, tt.want)
			}
		})
	}
}
