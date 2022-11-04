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

var msgStopServe = style.Faint.Render("Press the 'q' key to stop serve")

func initialChainServeModel(cmd *cobra.Command, session *cliui.Session) chainServeModel {
	bus := session.EventBus()
	ctx, quit := context.WithCancel(cmd.Context())

	// Update the command context to allow stopping by using the 'q' key
	cmd.SetContext(ctx)

	return chainServeModel{
		starting: true,
		status:   model.NewStatusEvents(bus, maxStatusEvents),
		events:   model.NewEvents(bus),
		cmd:      cmd,
		session:  session,
		quit:     quit,
	}
}

type chainServeModel struct {
	starting bool
	quitting bool
	broken   bool
	error    error
	cmd      *cobra.Command
	session  *cliui.Session
	quit     context.CancelFunc

	// Model definitions for the views
	status model.StatusEvents
	events model.Events
}

func (m chainServeModel) Init() tea.Cmd {
	// On initialization wait for status events and start serving the blockchain
	return tea.Batch(m.status.WaitEvent, chainServeStartCmd(m.cmd, m.session))
}

func (m chainServeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			m.quitting = true
			m.quit()

			// Remove the list of events received until now
			m.events.Clear()
		}
	case model.QuitMsg:
		cmd = tea.Quit
	case model.ErrorMsg:
		m.error = msg.Error
		cmd = tea.Quit
	case model.EventMsg:
		// TODO: See how to deal with code refresh when no errors

		// The first "starting" view displays status events one after the other
		// until the finish event is received, which signals that the second
		// "running" view must be displayed.
		if m.starting && msg.ProgressIndication == events.IndicationFinish {
			// Replace the starting view by the running one
			m.starting = false
			// Start waiting for events to display in the running view
			m.events, cmd = m.events.Update(msg)

			return m, cmd
		}

		if m.isRunning() {
			// Serve will keep running when there is an error after
			// a successful initialization until the code is fixed.
			if msg.Group == events.GroupError {
				m.broken = true

				// Clear the current events to only display the error
				m.events.Clear()
			} else if m.broken {
				m.broken = false

				// Clear the error event once the problem is fixed
				m.events.Clear()
			}
		}

		// Update the model that is being displayed
		return m.updateCurrentModel(msg)
	default:
		// Update the spinner of the model being displayed
		return m.updateCurrentModel(msg)
	}

	return m, cmd
}

func (m chainServeModel) View() string {
	// TODO: Generalize error and quit behaviours
	if m.error != nil {
		return fmt.Sprintf("%s %s\n", icons.NotOK, colors.Error(m.error.Error()))
	}

	if m.starting {
		return m.renderStartingView()
	}

	return m.renderRunningView()
}

func (m chainServeModel) isRunning() bool {
	return !(m.starting || m.quitting)
}

func (m chainServeModel) updateCurrentModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update the model that is being displayed
	if m.starting {
		m.status, cmd = m.status.Update(msg)
	} else {
		m.events, cmd = m.events.Update(msg)
	}

	return m, cmd
}

func (m chainServeModel) renderStartingView() string {
	var view strings.Builder

	view.WriteString(m.status.View())
	fmt.Fprintf(&view, "\n%s\n", msgStopServe)

	return model.FormatView(view.String())
}

func (m chainServeModel) renderRunningView() string {
	var view strings.Builder

	if m.quitting {
		if s := m.events.View(); s != "" {
			view.WriteString(s)
			view.WriteRune(model.EOL)
		}

		// TODO: Replace colors by lipgloss styles
		fmt.Fprintf(&view, "%s %s\n", icons.Info, colors.Info("Stopped"))
	} else {
		if !m.broken {
			view.WriteString("Chain is running\n\n")
		}

		view.WriteString(m.events.View())

		if m.broken {
			view.WriteString(colors.Info("\nWaiting for a fix before retrying...\n"))
		}

		fmt.Fprintf(&view, "\n%s\n", msgStopServe)
	}

	return model.FormatView(view.String())
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
