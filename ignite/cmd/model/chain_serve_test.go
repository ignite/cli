package cmdmodel_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"

	cmdmodel "github.com/ignite/cli/ignite/cmd/model"
	"github.com/ignite/cli/ignite/cmd/model/testdata"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	cliuimodel "github.com/ignite/cli/ignite/pkg/cliui/model"
	"github.com/ignite/cli/ignite/pkg/events"
)

var chainServeActions = colors.Faint("Press the 'q' key to stop serve")

func TestChainServeErrorView(t *testing.T) {
	// Arrange
	var model tea.Model

	err := errors.New("Test error")
	model = cmdmodel.NewChainServe(testdata.ModelContext{}, testdata.DummyEventsProvider{}, testdata.FooCmd)
	want := fmt.Sprintf("%s %s\n", icons.NotOK, colors.Error(err.Error()))

	// Arrange: Update model with an error message
	model, _ = model.Update(cliuimodel.ErrorMsg{Error: err})

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}

func TestChainServeStartView(t *testing.T) {
	// Arrange
	var model tea.Model

	spinner := cliuimodel.NewSpinner()
	queue := []string{"Event 1...", "Event 2..."}
	model = cmdmodel.NewChainServe(testdata.ModelContext{}, testdata.DummyEventsProvider{}, testdata.FooCmd)

	want := fmt.Sprintf(
		"%s%s\n\n%s %s %s\n\n%s\n",
		spinner.View(),
		queue[1],
		icons.OK,
		strings.TrimSuffix(queue[0], "..."),
		colors.Faint("0s"),
		chainServeActions,
	)
	want = cliuimodel.FormatView(want)

	// Arrange: Update model with status events
	for _, s := range queue {
		model, _ = model.Update(cliuimodel.EventMsg{
			Event: events.New(s, events.ProgressStart()),
			Start: time.Now(),
		})
	}

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}

func TestChainServeRunView(t *testing.T) {
	// Arrange
	var model tea.Model

	queue := []string{"Event 1", "Event 2"}
	model = cmdmodel.NewChainServe(testdata.ModelContext{}, testdata.DummyEventsProvider{}, testdata.FooCmd)

	want := fmt.Sprintf("Blockchain is running\n\n%s\n%s\n\n%s\n", queue[0], queue[1], chainServeActions)
	want = cliuimodel.FormatView(want)

	// Arrange: Update model with events
	for _, s := range queue {
		model, _ = model.Update(cliuimodel.EventMsg{
			Event: events.New(s, events.ProgressFinish()),
		})
	}

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}

func TestChainServeRunBrokenView(t *testing.T) {
	// Arrange
	var model tea.Model

	model = cmdmodel.NewChainServe(testdata.ModelContext{}, testdata.DummyEventsProvider{}, testdata.FooCmd)
	traceback := "Error traceback\nFoo"
	waitingFix := colors.Info("Waiting for a fix before retrying...")

	want := fmt.Sprintf("%s\n\n%s\n\n%s\n", traceback, waitingFix, chainServeActions)
	want = cliuimodel.FormatView(want)

	// Arrange: Update model to display the run view
	model, _ = model.Update(cliuimodel.EventMsg{
		Event: events.New("Run", events.ProgressFinish()),
	})

	// Arrange: Update model to display traceback within the run view
	model, _ = model.Update(cliuimodel.EventMsg{
		Event: events.New(traceback, events.Group(events.GroupError)),
	})

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}

func TestChainServeRebuildView(t *testing.T) {
	// Arrange
	var model tea.Model

	spinner := cliuimodel.NewSpinner()
	duration := colors.Faint("0s")
	queue := []string{"Event 1", "Event 2"}
	model = cmdmodel.NewChainServe(testdata.ModelContext{}, testdata.DummyEventsProvider{}, testdata.FooCmd)

	want := fmt.Sprintf(
		"Changes detected, restarting...\n\n%s%s\n\n%s %s %s\n%s Rebuild %s\n\n%s\n",
		spinner.View(),
		queue[1],
		icons.OK,
		queue[0],
		duration,
		icons.OK,
		duration,
		chainServeActions,
	)
	want = cliuimodel.FormatView(want)

	// Arrange: Update model to display the run view
	model, _ = model.Update(cliuimodel.EventMsg{
		Event: events.New("Run", events.ProgressFinish()),
	})

	// Arrange: Update model to display the rebuild view
	model, _ = model.Update(cliuimodel.EventMsg{
		Event: events.New("Rebuild", events.ProgressStart()),
		Start: time.Now(),
	})

	// Arrange: Update model with a status events
	for _, s := range queue {
		model, _ = model.Update(cliuimodel.EventMsg{
			Event: events.New(s, events.ProgressUpdate()),
			Start: time.Now(),
		})
	}

	// Act
	view := model.View()

	// Assert
	require.Equal(t, want, view)
}
