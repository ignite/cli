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
	"github.com/ignite/cli/ignite/pkg/xstrings"
)

const (
	stateChainDebugStarting uint = iota
	stateChainDebugRunning
)

var msgStopDebug = colors.Faint("Press the 'q' key to stop debug server")

// NewChainDebug returns a new UI model for the chain debug command.
func NewChainDebug(mCtx Context, bus events.Provider, cmd tea.Cmd) ChainDebug {
	// Initialize a context and cancel function to stop execution
	ctx, quit := context.WithCancel(mCtx.Context())

	// Update the context to allow stopping by using the 'q' key
	mCtx.SetContext(ctx)

	return ChainDebug{
		cmd:   cmd,
		quit:  quit,
		model: cliuimodel.NewEvents(bus),
	}
}

// ChainDebug defines a UI model for the chain debug command.
type ChainDebug struct {
	cmd   tea.Cmd
	quit  context.CancelFunc
	state uint
	error error
	model cliuimodel.Events
}

// Init is the first function that will be called.
// It returns a batch command that listen events and also runs the debug server.
func (m ChainDebug) Init() tea.Cmd {
	return tea.Batch(m.model.WaitEvent, m.cmd)
}

// Update is called when a message is received.
// It handles messages and executes the logic that updates the model.
func (m ChainDebug) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cliuimodel.QuitMsg:
		return m.processQuitMsg(msg)
	case cliuimodel.ErrorMsg:
		return m.processErrorMsg(msg)
	case tea.KeyMsg:
		return m.processKeyMsg(msg)
	case cliuimodel.EventMsg:
		return m.processEventMsg(msg)
	}

	return m.updateModel(msg)
}

// View renders the UI after every update.
func (m ChainDebug) View() string {
	if m.error != nil {
		s := xstrings.ToUpperFirst(m.error.Error())
		return fmt.Sprintf("%s %s\n", icons.NotOK, colors.Error(s))
	}

	var view strings.Builder

	switch m.state {
	case stateChainServeStarting:
		view.WriteString(m.renderStartView())
	case stateChainServeRunning:
		view.WriteString(m.renderRunView())
		view.WriteString(m.renderActions())
	}

	return cliuimodel.FormatView(view.String())
}

func (m ChainDebug) updateModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m ChainDebug) processQuitMsg(msg cliuimodel.QuitMsg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m ChainDebug) processErrorMsg(msg cliuimodel.ErrorMsg) (tea.Model, tea.Cmd) {
	m.error = msg.Error
	return m, tea.Quit
}

func (m ChainDebug) processKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if checkQuitKeyMsg(msg) {
		m.quit()
	}

	return m, nil
}

func (m ChainDebug) processEventMsg(msg cliuimodel.EventMsg) (tea.Model, tea.Cmd) {
	if m.state == stateChainDebugStarting {
		// Start view displays status events until the debug server is running.
		// When the status finish event is not an error it means that the debug
		// server started successfully and the run view is displayed.
		if msg.ProgressIndication == events.IndicationFinish {
			m.model.ClearEvents()
			m.state = stateChainDebugRunning
		}
	}

	return m.updateModel(msg)
}

func (m ChainDebug) renderActions() string {
	return fmt.Sprintf("\n%s\n", msgStopDebug)
}

func (m ChainDebug) renderStartView() string {
	return m.model.View()
}

func (m ChainDebug) renderRunView() string {
	var view strings.Builder

	view.WriteString("Blockchain is running\n\n")
	view.WriteString(m.model.View())

	return view.String()
}
