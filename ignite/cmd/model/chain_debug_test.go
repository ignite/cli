package cmdmodel_test

import (
	"errors"
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"

	cmdmodel "github.com/ignite/cli/ignite/cmd/model"
	"github.com/ignite/cli/ignite/cmd/model/testdata"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	cliuimodel "github.com/ignite/cli/ignite/pkg/cliui/model"
	"github.com/ignite/cli/ignite/pkg/events"
)

func TestChainDebugErrorView(t *testing.T) {
	// Arrange
	var model tea.Model

	err := errors.New("Test error")
	model = cmdmodel.NewChainDebug(testdata.ModelContext{}, testdata.DummyEventsProvider{}, testdata.FooCmd)
	want := fmt.Sprintf("%s %s\n", icons.NotOK, colors.Error(err.Error()))

	// Arrange: Update model with an error message
	model, _ = model.Update(cliuimodel.ErrorMsg{Error: err})

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}

func TestChainDebugStartView(t *testing.T) {
	// Arrange
	var model tea.Model

	spinner := cliuimodel.NewSpinner()
	queue := []string{"Event 1...", "Event 2..."}
	model = cmdmodel.NewChainDebug(testdata.ModelContext{}, testdata.DummyEventsProvider{}, testdata.FooCmd)

	want := fmt.Sprintf("\n%s%s\n", spinner.View(), queue[1])
	want = cliuimodel.FormatView(want)

	// Arrange: Update model with status events
	for _, s := range queue {
		model, _ = model.Update(cliuimodel.EventMsg{
			Event: events.New(s, events.ProgressStart()),
		})
	}

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}

func TestChainDebugRunView(t *testing.T) {
	// Arrange
	var model tea.Model

	evt := "Debug server: tcp://127.0.0.1:30500"
	actions := colors.Faint("Press the 'q' key to stop debug server")
	model = cmdmodel.NewChainDebug(testdata.ModelContext{}, testdata.DummyEventsProvider{}, testdata.FooCmd)

	want := fmt.Sprintf("Blockchain is running\n\n%s\n\n%s\n", evt, actions)
	want = cliuimodel.FormatView(want)

	// Arrange: Update model with a server running event
	model, _ = model.Update(cliuimodel.EventMsg{
		Event: events.New(evt, events.ProgressFinish()),
	})

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}
