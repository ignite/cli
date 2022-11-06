package ignitecmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cliui/model"
	"github.com/ignite/cli/ignite/pkg/cliui/style"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/services/chain"
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
	msgWaitingFix = colors.Info("\nWaiting for a fix before retrying...\n") // TODO: Replace colors by lipgloss styles
)

func initialChainServeModel(cmd *cobra.Command, session *cliui.Session) chainServeModel {
	bus := session.EventBus()
	ctx, quit := context.WithCancel(cmd.Context())

	// Update the command context to allow stopping by using the 'q' key
	cmd.SetContext(ctx)

	return chainServeModel{
		cmd:          cmd,
		session:      session,
		quit:         quit,
		startModel:   model.NewStatusEvents(bus, maxStatusEvents),
		runModel:     model.NewEvents(bus),
		rebuildModel: model.NewStatusEvents(bus, maxStatusEvents),
		quitModel:    model.NewEvents(bus),
	}
}

type chainServeModel struct {
	state   uint // TODO: Use a state machine for the state workflow?
	broken  bool
	error   error
	cmd     *cobra.Command
	session *cliui.Session
	quit    context.CancelFunc

	// Model definitions for the views
	startModel   model.StatusEvents
	runModel     model.Events
	rebuildModel model.StatusEvents
	quitModel    model.Events
}

func (m chainServeModel) Init() tea.Cmd {
	// On initialization wait for status events and start serving the blockchain
	return tea.Batch(m.startModel.WaitEvent, chainServeStartCmd(m.cmd, m.session))
}

func (m chainServeModel) Update(msg tea.Msg) (mod tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case model.QuitMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		return m.keyMsgUpdate(msg), nil
	case model.ErrorMsg:
		return m.errorMsgUpdate(msg)
	case model.EventMsg:
		return m.eventMsgUpdate(msg)
	}

	// By default update the spinner of the model being displayed
	return m.updateCurrentModel(msg)
}

func (m chainServeModel) View() string {
	// TODO: Generalize error and quit behaviours to be reused in other models
	if m.error != nil {
		return fmt.Sprintf("%s %s\n", icons.NotOK, colors.Error(m.error.Error()))
	}

	var view string

	switch m.state {
	case stateChainServeStarting:
		view = m.renderStartView()
	case stateChainServeRunning:
		view = m.renderRunView()
	case stateChainServeRebuilding:
		view = m.renderRebuildView()
	case stateChainServeQuitting:
		view = m.renderQuitView()
	}

	return model.FormatView(view)
}

func (m *chainServeModel) setState(state uint) {
	m.state = state
}

func (m *chainServeModel) errorMsgUpdate(msg model.ErrorMsg) (tea.Model, tea.Cmd) {
	m.error = msg.Error

	return m, tea.Quit
}

func (m *chainServeModel) keyMsgUpdate(msg tea.KeyMsg) tea.Model {
	if k := msg.String(); k == "q" || k == "ctrl+c" {
		m.setState(stateChainServeQuitting)
		m.quit()
	}

	return m
}

func (m *chainServeModel) eventMsgUpdate(msg model.EventMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	isFinishEvent := msg.ProgressIndication == events.IndicationFinish
	isErrorEvent := msg.Group == events.GroupError

	// Blockchain start/build finished
	if m.state == stateChainServeStarting && isFinishEvent {
		m.setState(stateChainServeRunning)

		// Start waiting for events to display in the running view
		m.runModel, cmd = m.runModel.Update(msg)

		return m, cmd
	}

	// Blockchain rebuild finished
	if m.state == stateChainServeRebuilding {
		if isErrorEvent {
			m.broken = true
		}

		if isFinishEvent {
			m.setState(stateChainServeRunning)

			// Make sure there are not events from the previous render
			m.runModel.Clear()

			// Start waiting for events to display in the run view
			m.runModel, cmd = m.runModel.Update(msg)

			return m, cmd
		}
	}

	if m.state == stateChainServeRunning {
		// Serve will keep running when there is an error after
		// a successful initialization until the code is fixed.
		if isErrorEvent {
			// If an error event is received it means there is an issue with the source code
			m.broken = true

			// Clear the current events to only display the error
			m.runModel.Clear()

			m.runModel, cmd = m.runModel.Update(msg)

			return m, cmd
		}

		// When a status event is received during run it means something
		// changed in the source code to trigger the blockchain rebuild.
		if msg.InProgress() {
			m.setState(stateChainServeRebuilding)
			m.rebuildModel.Clear()

			m.rebuildModel, cmd = m.rebuildModel.Update(msg)

			return m, cmd
		}

		// When the source code is not working and the event
		// is not an error event it means the issue was fixed.
		if m.broken {
			m.broken = false

			// Clear the error events
			m.runModel.Clear()
		}
	}

	// Update the model that is being displayed
	return m.updateCurrentModel(msg)
}

func (m chainServeModel) updateCurrentModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update the model that is being displayed
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

func (m chainServeModel) renderActionsMenu() string {
	return fmt.Sprintf("\n%s\n", msgStopServe)
}

func (m chainServeModel) renderStartView() string {
	var view strings.Builder

	view.WriteString(m.startModel.View())
	view.WriteString(m.renderActionsMenu())

	return view.String()
}

func (m chainServeModel) renderRunView() string {
	var view strings.Builder

	if !m.broken {
		view.WriteString("Blockchain is running\n\n")
	}

	view.WriteString(m.runModel.View())

	if m.broken {
		view.WriteString(msgWaitingFix)
	}

	view.WriteString(m.renderActionsMenu())

	return view.String()
}

func (m chainServeModel) renderRebuildView() string {
	var view strings.Builder

	if !m.broken {
		view.WriteString("Changes detected, restarting...\n\n")
	}

	view.WriteString(m.rebuildModel.View())

	if m.broken {
		view.WriteString(msgWaitingFix)
	}

	view.WriteString(m.renderActionsMenu())

	return view.String()
}

func (m chainServeModel) renderQuitView() string {
	var view strings.Builder

	// Display the events received during quit
	if s := m.quitModel.View(); s != "" {
		view.WriteString(s)
		view.WriteRune(model.EOL)
	}

	fmt.Fprintf(&view, "%s %s\n", icons.Info, colors.Info("Stopped"))

	return view.String()
}

func chainServeStartCmd(cmd *cobra.Command, session *cliui.Session) tea.Cmd {
	return func() tea.Msg {
		chainOption := []chain.Option{
			chain.WithOutputer(session),
			chain.CollectEvents(session.EventBus()),
		}

		if flagGetProto3rdParty(cmd) {
			chainOption = append(chainOption, chain.EnableThirdPartyModuleCodegen())
		}

		if flagGetCheckDependencies(cmd) {
			chainOption = append(chainOption, chain.CheckDependencies())
		}

		// check if custom config is defined
		config, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			return err
		}
		if config != "" {
			chainOption = append(chainOption, chain.ConfigFile(config))
		}

		// create the chain
		c, err := NewChainWithHomeFlags(cmd, chainOption...)
		if err != nil {
			return err
		}

		cacheStorage, err := newCache(cmd)
		if err != nil {
			return err
		}

		// serve the chain
		var serveOptions []chain.ServeOption

		forceUpdate, err := cmd.Flags().GetBool(flagForceReset)
		if err != nil {
			return err
		}

		if forceUpdate {
			serveOptions = append(serveOptions, chain.ServeForceReset())
		}

		resetOnce, err := cmd.Flags().GetBool(flagResetOnce)
		if err != nil {
			return err
		}

		if resetOnce {
			serveOptions = append(serveOptions, chain.ServeResetOnce())
		}

		quitOnFail, err := cmd.Flags().GetBool(flagQuitOnFail)
		if err != nil {
			return err
		}

		if quitOnFail {
			serveOptions = append(serveOptions, chain.QuitOnFail())
		}

		if flagGetSkipProto(cmd) {
			serveOptions = append(serveOptions, chain.ServeSkipProto())
		}

		err = c.Serve(cmd.Context(), cacheStorage, serveOptions...)
		if err != nil && !errors.Is(err, context.Canceled) {
			return model.ErrorMsg{Error: err}
		}

		return model.QuitMsg{}
	}
}
