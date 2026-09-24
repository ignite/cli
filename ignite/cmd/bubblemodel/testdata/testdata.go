package testdata

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/ignite/cli/v30/ignite/pkg/events"
)

func FooCmd() tea.Msg { return nil }

type ModelContext struct{}

func (ModelContext) Context() context.Context   { return context.TODO() }
func (ModelContext) SetContext(context.Context) {}

type DummyEventsProvider struct{}

func (DummyEventsProvider) Events() <-chan events.Event {
	c := make(chan events.Event)
	close(c)
	return c
}
