package cmdmodel

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	cliuimodel "github.com/ignite/cli/ignite/pkg/cliui/model"
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
	msgStopServe  = colors.Faint("Press the 'q' key to stop serve")
	msgWaitingFix = colors.Info("Waiting for a fix before retrying...")
)

type Context interface {
	// Context returns the current context.
	Context() context.Context

	// SetContext updates the context with a new one.
	SetContext(context.Context)
}

// NewChainServe returns a new UI model for the chain serve command.
func NewChainServe(mCtx Context, bus events.Provider, cmd tea.Cmd) ChainServe {
	// Initialize a context and cancel function to stop execution
	ctx, quit := context.WithCancel(mCtx.Context())

	// Update the context to allow stopping by using the 'q' key
	mCtx.SetContext(ctx)

	return ChainServe{
		cmd:          cmd,
		quit:         quit,
		startModel:   cliuimodel.NewStatusEvents(bus, maxStatusEvents),
		runModel:     cliuimodel.NewEvents(bus),
		rebuildModel: cliuimodel.NewStatusEvents(bus, maxStatusEvents),
		quitModel:    cliuimodel.NewEvents(bus),
	}
}

// ChainServe defines a UI model for the chain serve command.
type ChainServe struct {
	cmd  tea.Cmd
	quit context.CancelFunc

	state  uint  // Keeps track of the model/view being displayed
	broken bool  // True when blockchain app's source code has issues
	error  error // Critical error returned during command execution

	// Model definitions for the chain serve views
	startModel   cliuimodel.StatusEvents
	runModel     cliuimodel.Events
	rebuildModel cliuimodel.StatusEvents
	quitModel    cliuimodel.Events
}

// Init is the first function that will be called.
// It returns a batch command that listen events and also runs the blockchain app.
func (m ChainServe) Init() tea.Cmd {
	// On initialization wait for status events and start serving the blockchain
	return tea.Batch(m.startModel.WaitEvent, m.cmd)
}

// Update is called when a message is received.
// It handles messages and executes the logic that updates the model.
func (m ChainServe) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if checkQuitKeyMsg(msg) {
		m.state = stateChainServeQuitting
	}

	switch msg := msg.(type) {
	case cliuimodel.QuitMsg:
		return m.processQuitMsg(msg)
	case cliuimodel.ErrorMsg:
		return m.processErrorMsg(msg)
	case tea.KeyMsg:
		return m.processKeyMsg(msg)
	case cliuimodel.EventMsg:
		return m.processEventMsg(msg)
	default:
		return m.updateCurrentModel(msg)
	}
}

// View renders the UI after every update.
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
		view.WriteString(m.renderActions())
	}

	return cliuimodel.FormatView(view.String())
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

func (m ChainServe) processQuitMsg(cliuimodel.QuitMsg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m ChainServe) processErrorMsg(msg cliuimodel.ErrorMsg) (tea.Model, tea.Cmd) {
	m.error = msg.Error
	return m, tea.Quit
}

func (m ChainServe) processKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if checkQuitKeyMsg(msg) {
		// Cancel the context to signal stop
		m.quit()
	}

	return m, nil
}

func (m ChainServe) processEventMsg(msg cliuimodel.EventMsg) (tea.Model, tea.Cmd) {
	// When an error event is received it means there is an issue with
	// the blockchain app's source code that the user must fix.
	m.broken = msg.Group == events.GroupError

	// UI responds to key press or mouse events by default but we use
	// events and the events bus to interact with the UI during execution.
	// Check if the state must be changed to switch to a different view.
	switch m.state {
	case stateChainServeStarting:
		// Start view displays status events until the blockchain is running or an
		// error event is received in which case it displays the run view with an
		// error traceback and waits until the issue is fixed.
		// When the status finish event is not an error it means that the blockchain
		// started successfully and the run view is displayed.
		if msg.ProgressIndication == events.IndicationFinish {
			m.state = stateChainServeRunning
		}
	case stateChainServeRunning:
		// Run view shows account addresses, API URLs and the paths required to
		// have a context on the running blockchain app and waits for errors or
		// changes in the blockchain app source code.
		// If an error event is received during run it means that there is an error
		// in the app source code in which case the error message and traceback are
		// displayed until the code is fixed, or otherwise when an status event is
		// received it means that the code changed so the app must be rebuilt.
		if m.broken {
			// Clear events to only display the error received with the last event message
			m.runModel.ClearEvents()
		} else if msg.InProgress() {
			// When a status event is received during run it means something
			// changed in the source code which triggers the blockchain rebuild.
			m.runModel.ClearEvents()
			m.state = stateChainServeRebuilding
		}
	case stateChainServeRebuilding:
		// Rebuild view is similar to run view but only displayed when the source
		// code changes and the blockchain is rebuilt.
		// When the status finish event is not an error it means that the blockchain
		// was rebuilt successfully and the run view is displayed.
		if msg.ProgressIndication == events.IndicationFinish {
			m.rebuildModel.ClearEvents()
			m.state = stateChainServeRunning
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

func checkQuitKeyMsg(m tea.Msg) bool {
	msg, ok := m.(tea.KeyMsg)
	if !ok {
		return false
	}

	key := msg.String()

	return key == "q" || key == "ctrl+c"
}
