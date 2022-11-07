package cmdmodel

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	cliuimodel "github.com/ignite/cli/ignite/pkg/cliui/model"
)

type Context interface {
	// Context returns the current context.
	Context() context.Context

	// SetContext updates the context with a new one.
	SetContext(context.Context)
}

func newModel(mCtx Context, cmd tea.Cmd) model {
	// Initialize a context and cancel function to stop execution
	ctx, quit := context.WithCancel(mCtx.Context())

	// Update the context to allow stopping by using the 'q' key
	mCtx.SetContext(ctx)

	return model{
		cmd:  cmd,
		quit: quit,
	}
}

type model struct {
	cmd  tea.Cmd
	quit context.CancelFunc
}

func (m model) Init() tea.Cmd {
	return m.cmd
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if checkQuitKeyMsg(msg) {
			// Cancel the context to signal stop
			m.quit()
		}
	case cliuimodel.QuitMsg:
		// When a quit message is received stop execution
		return m, tea.Quit
	}

	return m, nil
}

func checkQuitKeyMsg(m tea.Msg) bool {
	msg, ok := m.(tea.KeyMsg)
	if !ok {
		return false
	}

	key := msg.String()

	return key == "q" || key == "ctrl+c"
}
