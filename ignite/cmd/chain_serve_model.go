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
	error    error
	status   model.StatusEvents
	events   model.Events
	cmd      *cobra.Command
	session  *cliui.Session
	quit     context.CancelFunc
}

func (m chainServeModel) Init() tea.Cmd {
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
		// TODO: See how this works with the running view
		// TODO: How to deal with non static events like the temporary serve restart error ones
		if msg.ProgressIndication == events.IndicationFinish {
			// Replace the starting view by the running one
			m.starting = false
			// Start waiting for events to display in the running view
			m.events, cmd = m.events.Update(msg)

			return m, cmd
		}

		if m.starting {
			m.status, cmd = m.status.Update(msg)
		} else {
			m.events, cmd = m.events.Update(msg)
		}
	default:
		// This is required to allow event spinner updates
		m.status, cmd = m.status.Update(msg)
	}

	return m, cmd
}

func (m chainServeModel) View() string {
	// TODO: Generalize the error and quit behaviour
	if m.error != nil {
		return fmt.Sprintf("%s %s\n", icons.NotOK, colors.Error(m.error.Error()))
	}

	if m.starting {
		return m.renderStartingView()
	}

	return m.renderRunningView()
}

func (m chainServeModel) renderStartingView() string {
	var view strings.Builder

	view.WriteString(m.status.View())
	fmt.Fprintf(&view, "%s\n", msgStopServe)

	return model.FormatView(view.String())
}

func (m chainServeModel) renderRunningView() string {
	var view strings.Builder

	if m.quitting {
		// TODO: Add spinner to the models.Events for ongoing events
		view.WriteString(m.events.View())
		// TODO: Replace colors by lipgloss styles?
		fmt.Fprintf(&view, "%s %s\n", icons.Info, colors.Info("Stopped"))
	} else {
		view.WriteString("Chain is running\n\n")
		view.WriteString(m.events.View())
		fmt.Fprintf(&view, "%s\n", msgStopServe)
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
		c, err := newChainWithHomeFlags(cmd, chainOption...)
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
