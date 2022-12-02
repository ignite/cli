package cliuimodel_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	cliuimodel "github.com/ignite/cli/ignite/pkg/cliui/model"
	"github.com/ignite/cli/ignite/pkg/events"
)

func TestStatusEventsView(t *testing.T) {
	// Arrange
	spinner := cliuimodel.NewSpinner()
	queue := []string{"Event 1...", "Event 2..."}
	model := cliuimodel.NewStatusEvents(dummyEventsProvider{}, len(queue))
	want := fmt.Sprintf(
		"Static event\n\n%s%s\n\n%s %s %s\n",
		spinner.View(),
		queue[1],
		icons.OK,
		strings.TrimSuffix(queue[0], "..."),
		colors.Faint("0s"),
	)

	// Arrange: Update model with status events
	for _, s := range queue {
		model, _ = model.Update(cliuimodel.EventMsg{
			Event: events.New(s, events.ProgressStart()),
			Start: time.Now(),
		})
	}

	// Arrange: Add one static event
	model, _ = model.Update(cliuimodel.EventMsg{
		Event: events.New("Static event"),
	})

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}

func TestEventsView(t *testing.T) {
	// Arrange
	spinner := cliuimodel.NewSpinner()
	model := cliuimodel.NewEvents(dummyEventsProvider{})
	queue := []string{"Event 1", "Event 2"}
	want := fmt.Sprintf(
		"%s\n%s\n\n%sStatus Event...\n",
		queue[0],
		queue[1],
		spinner.View(),
	)

	// Arrange: Update model with events
	for _, s := range queue {
		model, _ = model.Update(cliuimodel.EventMsg{
			Event: events.New(s),
		})
	}

	// Arrange: Add one status event
	model, _ = model.Update(cliuimodel.EventMsg{
		Event: events.New("Status Event...", events.ProgressStart()),
	})

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}

type dummyEventsProvider struct{}

func (dummyEventsProvider) Events() <-chan events.Event {
	c := make(chan events.Event)
	close(c)
	return c
}
