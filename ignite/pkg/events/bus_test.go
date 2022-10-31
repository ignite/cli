package events_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/events"
)

func TestBusSend(t *testing.T) {
	cases := []struct {
		name, message string
		options       []events.Option
		progress      events.ProgressIndication
	}{
		{
			name:     "without options",
			message:  "test",
			progress: events.IndicationNone,
		},
		{
			name:     "with options",
			message:  "test",
			options:  []events.Option{events.ProgressStart()},
			progress: events.IndicationStart,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			bus := events.NewBus()
			defer bus.Stop()

			// Act
			bus.Send(tt.message, tt.options...)

			// Assert
			select {
			case e := <-bus.Events():
				require.Equal(t, tt.message, e.Message)
				require.Equal(t, tt.progress, e.ProgressIndication)
			default:
				t.Error("expected an event to be received")
			}
		})
	}
}

func TestBusSendf(t *testing.T) {
	// Arrange
	bus := events.NewBus()
	defer bus.Stop()

	want := "foo 42"

	// Act
	bus.Sendf("%s %d", "foo", 42)

	// Assert
	select {
	case e := <-bus.Events():
		require.Equal(t, want, e.Message)
	default:
		t.Error("expected an event to be received")
	}
}

func TestBusSendInfo(t *testing.T) {
	cases := []struct {
		name, message string
		options       []events.Option
		progress      events.ProgressIndication
	}{
		{
			name:     "without options",
			message:  "test",
			progress: events.IndicationNone,
		},
		{
			name:     "with options",
			message:  "test",
			options:  []events.Option{events.ProgressStart()},
			progress: events.IndicationStart,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			bus := events.NewBus()
			defer bus.Stop()

			// Act
			bus.SendInfo(tt.message, tt.options...)

			// Assert
			select {
			case e := <-bus.Events():
				require.Equal(t, colors.Info(tt.message), e.Message)
				require.Equal(t, tt.progress, e.ProgressIndication)
			default:
				t.Error("expected an event to be received")
			}
		})
	}
}

func TestBusSendError(t *testing.T) {
	cases := []struct {
		name, message string
		options       []events.Option
		progress      events.ProgressIndication
	}{
		{
			name:     "without options",
			message:  "test",
			progress: events.IndicationNone,
		},
		{
			name:     "with options",
			message:  "test",
			options:  []events.Option{events.ProgressStart()},
			progress: events.IndicationStart,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			bus := events.NewBus()
			defer bus.Stop()

			err := errors.New(tt.message)

			// Act
			bus.SendError(err, tt.options...)

			// Assert
			select {
			case e := <-bus.Events():
				require.Equal(t, colors.Error(tt.message), e.Message)
				require.Equal(t, tt.progress, e.ProgressIndication)
			default:
				t.Error("expected an event to be received")
			}
		})
	}
}

type testEventView struct {
	message string
}

func (v testEventView) String() string {
	return v.message
}

func TestBusSendView(t *testing.T) {
	cases := []struct {
		name, message string
		options       []events.Option
		progress      events.ProgressIndication
	}{
		{
			name:     "without options",
			message:  "test",
			progress: events.IndicationNone,
		},
		{
			name:     "with options",
			message:  "test",
			options:  []events.Option{events.ProgressStart()},
			progress: events.IndicationStart,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			bus := events.NewBus()
			defer bus.Stop()

			view := testEventView{tt.message}

			// Act
			bus.SendView(view, tt.options...)

			// Assert
			select {
			case e := <-bus.Events():
				require.Equal(t, tt.message, e.Message)
				require.Equal(t, tt.progress, e.ProgressIndication)
			default:
				t.Error("expected an event to be received")
			}
		})
	}
}

func TestBusStop(t *testing.T) {
	// Arrange
	bus := events.NewBus()

	// Act
	bus.Stop()
	bus.Send("ignored message")
	_, ok := <-bus.Events()

	// Assert
	require.False(t, ok, "expected no events after bus stopped")
}
