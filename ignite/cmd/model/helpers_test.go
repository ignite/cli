package cmdmodel_test

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ignite/cli/ignite/pkg/events"
)

func fooCmd() tea.Msg { return nil }

type modelContext struct{}

func (modelContext) Context() context.Context   { return context.TODO() }
func (modelContext) SetContext(context.Context) {}

type dummyEventsProvider struct{}

func (dummyEventsProvider) Events() <-chan events.Event {
	c := make(chan events.Event)
	close(c)
	return c
}
