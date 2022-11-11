package cmdmodel

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	cliuimodel "github.com/ignite/cli/ignite/pkg/cliui/model"
	"github.com/ignite/cli/ignite/pkg/cliui/style"
	"github.com/ignite/cli/ignite/pkg/events"
)

const (
	maxStatusEvents = 7
)

const (
	stateChainServeStarting uint = iota
	stateChainServeRunning
	stateChainServeRebuilding
	stateChainServeQuitting
)

var (
	msgStopServe  = style.Faint.Render("Press the 'q' key to stop serve")
	msgWaitingFix = colors.Info("Waiting for a fix before retrying...") // TODO: Replace colors by lipgloss styles
)

// NewChainServe returns a new UI model for the chain serve command.
func NewChainServe(ctx Context, bus events.Provider, cmd tea.Cmd) ChainServe {
	return ChainServe{
		model:        newModel(ctx, cmd),
		startModel:   cliuimodel.NewStatusEvents(bus, maxStatusEvents),
		runModel:     cliuimodel.NewEvents(bus),
		rebuildModel: cliuimodel.NewStatusEvents(bus, maxStatusEvents),
		quitModel:    cliuimodel.NewEvents(bus),
	}
}

// ChainServe defines a UI model for the chain serve command.
type ChainServe struct {
	model

	state  uint  // Keeps track of the model/view being displayed
	broken bool  // True when blockchain app's source code has issues
	error  error // Critical error returned during command execution

	// Model definitions for the chain serve views
	startModel   cliuimodel.StatusEvents
	runModel     cliuimodel.Events
	rebuildModel cliuimodel.StatusEvents
	quitModel    cliuimodel.Events
}

func (m ChainServe) Init() tea.Cmd {
	// On initialization wait for status events and start serving the blockchain
	return tea.Batch(m.startModel.WaitEvent, m.model.Init())
}

func (m ChainServe) Update(msg tea.Msg) (mod tea.Model, cmd tea.Cmd) {
	if checkQuitKeyMsg(msg) {
		m.setState(stateChainServeQuitting)
	}

	switch msg := msg.(type) {
	case cliuimodel.EventMsg:
		return m.processEventMsg(msg)
	case cliuimodel.ErrorMsg:
		m.error = msg.Error
		return m, tea.Quit
	default:
		if m.model, cmd = m.model.Update(msg); cmd != nil {
			return m, cmd
		}

		// Update the model that is being displayed
		return m.updateCurrentModel(msg)
	}
}

func (m ChainServe) View() string {
	if m.error != nil {
		return fmt.Sprintf("%s %s\n", icons.NotOK, colors.Error(m.error.Error()))
	}

	var view strings.Builder

	switch m.state {
	case stateChainServeStarting:
		view.WriteString(m.renderStartView())
	case stateChainServeRunning:
		view.WriteString(m.renderRunView())
	case stateChainServeRebuilding:
		view.WriteString(m.renderRebuildView())
	case stateChainServeQuitting:
		view.WriteString(m.renderQuitView())
	}

	if m.state != stateChainServeQuitting {
		// TODO: Add actions to copy mnemonics to clipboard in run view?
		view.WriteString(m.renderActions())
	}

	return cliuimodel.FormatView(view.String())
}

func (m *ChainServe) setState(state uint) {
	m.state = state
}

func (m ChainServe) updateCurrentModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case stateChainServeStarting:
		m.startModel, cmd = m.startModel.Update(msg)
	case stateChainServeRunning:
		m.runModel, cmd = m.runModel.Update(msg)
	case stateChainServeRebuilding:
		m.rebuildModel, cmd = m.rebuildModel.Update(msg)
	case stateChainServeQuitting:
		m.quitModel, cmd = m.quitModel.Update(msg)
	}

	return m, cmd
}

func (m *ChainServe) processEventMsg(msg cliuimodel.EventMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// When an error event is received it means there is an issue with
	// the blockchain app's source code that the user must fix.
	m.broken = msg.Group == events.GroupError

	// UI responds to key press or mouse events by default but we use
	// events and the events bus to interact with the UI during execution.
	// State defines the model/view that is being displayed.
	switch m.state {
	case stateChainServeStarting:
		// Start view displays status events until the blockchain is running or an
		// error event is received in which case it displays the run view with an
		// error traceback and waits until the issue is fixed.
		// When the status finish event is not an error it means that the blockchain
		// started successfully and the run view is displayed.
		if msg.ProgressIndication == events.IndicationFinish {
			m.setState(stateChainServeRunning)

			// Prepare model for the run view
			m.runModel, cmd = m.runModel.Update(msg)

			return m, cmd
		}
	case stateChainServeRunning:
		// Serve will be displayed when there is an error until the code is fixed.
		// It will show an error traceback until the code is changed and the rebuild
		// view is displayed.
		if msg.Group == events.GroupError {
			// Clear the current events to only display the error
			m.runModel.ClearEvents()

			// Prepare model for the run view
			m.runModel, cmd = m.runModel.Update(msg)

			return m, cmd
		}

		// When the received event is not an error make sure to change the broken property
		if m.broken {
			m.broken = false

			// Clear error events
			m.runModel.ClearEvents()
		}

		// When a status event is received during run it means something
		// changed in the source code to trigger the blockchain rebuild.
		if msg.InProgress() {
			m.setState(stateChainServeRebuilding)

			// Make sure there are not events from a previous "rebuild" render
			m.rebuildModel.ClearEvents()

			// Prepare model for the rebuild view
			m.rebuildModel, cmd = m.rebuildModel.Update(msg)

			return m, cmd
		}
	case stateChainServeRebuilding:
		// Rebuild view is similar to run view but only displayed when the source
		// code changes and the blockchain is rebuilt.
		// When the status finish event is not an error it means that the blockchain
		// was rebuilt successfully and the run view is displayed.
		if msg.ProgressIndication == events.IndicationFinish {
			m.setState(stateChainServeRunning)

			// Make sure there are not events from a previous "run" render
			m.runModel.ClearEvents()

			// Prepare model for the run view
			m.runModel, cmd = m.runModel.Update(msg)

			return m, cmd
		}
	}

	// Update the model that is being displayed
	return m.updateCurrentModel(msg)
}

func (m ChainServe) renderActions() string {
	return fmt.Sprintf("\n%s\n", msgStopServe)
}

func (m ChainServe) renderStartView() string {
	return m.startModel.View()
}

func (m ChainServe) renderRunView() string {
	var view strings.Builder

	if !m.broken {
		view.WriteString("Blockchain is running\n\n")
	}

	view.WriteString(m.runModel.View())

	if m.broken {
		fmt.Fprintf(&view, "\n%s\n", msgWaitingFix)
	}

	return view.String()
}

func (m ChainServe) renderRebuildView() string {
	var view strings.Builder

	if !m.broken {
		view.WriteString("Changes detected, restarting...\n\n")
	}

	view.WriteString(m.rebuildModel.View())

	if m.broken {
		fmt.Fprintf(&view, "\n%s\n", msgWaitingFix)
	}

	return view.String()
}

func (m ChainServe) renderQuitView() string {
	var view strings.Builder

	// Display the events received during quit
	if s := m.quitModel.View(); s != "" {
		view.WriteString(s)
		view.WriteRune('\n')
	}

	fmt.Fprintf(&view, "%s %s\n", icons.Info, colors.Info("Stopped"))

	return view.String()
}
