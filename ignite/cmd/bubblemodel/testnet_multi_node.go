package cmdmodel

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/v29/ignite/services/chain"
)

// NodeStatus is an integer data type that represents the status of a node.
type NodeStatus int

const (
	// Stopped indicates that the node is currently stopped.
	Stopped NodeStatus = iota

	// Running indicates that the node is currently running.
	Running
)

// ui styling constants.
var (
	// base colors.
	activeColor    = lipgloss.Color("#1B7FCA") // bright blue
	subtleColor    = lipgloss.Color("#5C6A72") // dark gray
	textColor      = lipgloss.Color("#232326") // nearly black
	highlightColor = lipgloss.Color("#10B981") // green
	warningColor   = lipgloss.Color("#FF5436") // red
	focusedColor   = lipgloss.Color("#A27DF8") // purple

	// tabs styling.
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	tabStyle = lipgloss.NewStyle().
			Border(tabBorder).
			BorderForeground(subtleColor).
			Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Border(activeTabBorder).
			BorderForeground(activeColor).
			Foreground(activeColor).
			Bold(true).
			Padding(0, 1)

	// active/stopped tab styles.
	runningTabStyle = lipgloss.NewStyle().
			Border(tabBorder).
			BorderForeground(highlightColor).
			Foreground(subtleColor).
			Padding(0, 1)

	activeRunningTabStyle = lipgloss.NewStyle().
				Border(activeTabBorder).
				BorderForeground(highlightColor).
				Foreground(highlightColor).
				Bold(true).
				Padding(0, 1)

	// node status styles.
	nodeActiveStyle  = lipgloss.NewStyle().Foreground(highlightColor).Bold(true)
	nodeStoppedStyle = lipgloss.NewStyle().Foreground(warningColor)
	tcpStyle         = lipgloss.NewStyle().Foreground(activeColor)
	infoStyle        = lipgloss.NewStyle().Foreground(subtleColor)

	// header styling.
	headerStyle = lipgloss.NewStyle().
			Foreground(focusedColor).
			Bold(true).
			Padding(0, 0, 1, 0)

	// log styles.
	logEntryStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2)

	logBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor).
			Padding(1, 2).
			Width(80)
)

// Make sure MultiNode implements tea.Model interface.
var _ tea.Model = MultiNode{}

// MultiNode represents a set of nodes, managing state and information related to them.
type MultiNode struct {
	ctx  context.Context
	appd string
	args chain.MultiNodeArgs

	nodeStatuses []NodeStatus
	pids         []int      // Store the PIDs of the running processes
	numNodes     int        // Number of nodes
	logs         [][]string // Store logs for each node

	// UI state
	selectedNode int        // Currently selected node index
	help         help.Model // Help menu model
	showHelp     bool       // Whether to show the help menu
}

// ToggleNodeMsg is a structure used to pass messages
// to enable or disable a node based on the node index.
type ToggleNodeMsg struct {
	nodeIdx int
}

// UpdateStatusMsg defines a message that updates the status of a node by index.
type UpdateStatusMsg struct {
	nodeIdx int
	status  NodeStatus
}

// UpdateLogsMsg is for continuously updating the chain logs in the View.
type UpdateLogsMsg struct{}

// SwitchFocusMsg indicates a switch in focus to another node.
type SwitchFocusMsg struct {
	nodeIdx int
}

// UpdateDeemon returns a command that sends an UpdateLogsMsg.
// This command is intended to continuously refresh the logs displayed in the user interface.
func UpdateDeemon() tea.Cmd {
	return func() tea.Msg {
		return UpdateLogsMsg{}
	}
}

// NewModel initializes the model.
func NewModel(ctx context.Context, chainname string, args chain.MultiNodeArgs) (MultiNode, error) {
	numNodes, err := strconv.Atoi(args.NumValidator)
	if err != nil {
		return MultiNode{}, err
	}

	h := help.New()
	h.ShowAll = true

	return MultiNode{
		ctx:          ctx,
		appd:         chainname + "d",
		args:         args,
		nodeStatuses: make([]NodeStatus, numNodes), // initial states of nodes
		pids:         make([]int, numNodes),
		numNodes:     numNodes,
		logs:         make([][]string, numNodes), // Initialize logs for each node
		selectedNode: 0,                          // Select the first node initially
		help:         h,
		showHelp:     false,
	}, nil
}

// Init implements the Init method of the tea.Model interface.
func (m MultiNode) Init() tea.Cmd {
	// start all nodes as soon as the application launches
	return m.StartAllNodes()
}

// ToggleNode toggles the state of a node.
func ToggleNode(nodeIdx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleNodeMsg{nodeIdx: nodeIdx}
	}
}

// SwitchFocus changes the focus to a specific node.
func SwitchFocus(nodeIdx int) tea.Cmd {
	return func() tea.Msg {
		return SwitchFocusMsg{nodeIdx: nodeIdx}
	}
}

// RunNode runs or stops the node based on its status.
func RunNode(nodeIdx int, start bool, m MultiNode) tea.Cmd {
	var (
		pid  = &m.pids[nodeIdx]
		args = m.args
		appd = m.appd
	)

	return func() tea.Msg {
		if start {
			nodeHome := filepath.Join(args.OutputDir, args.NodeDirPrefix+strconv.Itoa(nodeIdx))
			// Create the command to run in the background as a daemon
			cmd := exec.Command(appd, "start", "--home", nodeHome)

			// Start the process as a daemon
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setpgid: true, // Ensure it runs in a new process group
			}

			stdout, err := cmd.StdoutPipe() // Get stdout for logging
			if err != nil {
				fmt.Printf("Failed to start node %d: %v\n", nodeIdx+1, err)
				return UpdateStatusMsg{nodeIdx: nodeIdx, status: Stopped}
			}

			err = cmd.Start() // Start the node in the background
			if err != nil {
				fmt.Printf("Failed to start node %d: %v\n", nodeIdx+1, err)
				return UpdateStatusMsg{nodeIdx: nodeIdx, status: Stopped}
			}

			*pid = cmd.Process.Pid // Store the PID

			// Create an errgroup with context
			g, gCtx := errgroup.WithContext(m.ctx)
			g.Go(func() error {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					select {
					case <-gCtx.Done():
						// Handle context cancellation
						return gCtx.Err()
					default:
						line := scanner.Text()
						// Add log line to the respective node's log slice
						m.logs[nodeIdx] = append(m.logs[nodeIdx], line)
						// Keep only the last 5 lines
						if len(m.logs[nodeIdx]) > 5 {
							m.logs[nodeIdx] = m.logs[nodeIdx][len(m.logs[nodeIdx])-5:]
						}
					}
				}
				if err := scanner.Err(); err != nil {
					return err
				}
				return nil
			})

			// Goroutine to handle stopping the node if context is canceled
			g.Go(func() error {
				<-gCtx.Done() // Wait for context to be canceled

				// Stop the daemon process if context is canceled
				if *pid != 0 {
					err := syscall.Kill(-*pid, syscall.SIGTERM) // Stop the daemon process
					if err != nil {
						fmt.Printf("Failed to stop node %d: %v\n", nodeIdx+1, err)
					} else {
						*pid = 0 // Reset PID after stopping
					}
				}

				return gCtx.Err()
			})

			return UpdateStatusMsg{nodeIdx: nodeIdx, status: Running}
		}
		// Use kill to stop the node process by PID
		if *pid != 0 {
			err := syscall.Kill(-*pid, syscall.SIGTERM) // Stop the daemon process
			if err != nil {
				fmt.Printf("Failed to stop node %d: %v\n", nodeIdx+1, err)
			} else {
				*pid = 0 // Reset PID after stopping
			}
		}
		return UpdateStatusMsg{nodeIdx: nodeIdx, status: Stopped}
	}
}

// StopAllNodes stops all nodes.
func (m *MultiNode) StopAllNodes() {
	for i := range m.numNodes {
		if m.nodeStatuses[i] == Running {
			RunNode(i, false, *m)() // Stop node
		}
	}
}

// StartAllNodes starts all nodes that are currently stopped.
func (m *MultiNode) StartAllNodes() tea.Cmd {
	cmds := make([]tea.Cmd, 0, m.numNodes)
	for i := range m.numNodes {
		if m.nodeStatuses[i] == Stopped {
			cmds = append(cmds, RunNode(i, true, *m))
		}
	}

	return tea.Batch(cmds...)
}

// Update handles messages and updates the model.
func (m MultiNode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.StopAllNodes() // Stop all nodes before quitting
			return m, tea.Quit
		case "h":
			// Toggle help screen
			m.showHelp = !m.showHelp
			return m, nil
		case "tab", "right":
			// Move selection to the next node
			m.selectedNode = (m.selectedNode + 1) % m.numNodes
			return m, nil
		case "shift+tab", "left":
			// Move selection to the previous node
			m.selectedNode = (m.selectedNode - 1 + m.numNodes) % m.numNodes
			return m, nil
		default:
			// Check for numbers from 1 to numNodes
			for i := 0; i < m.numNodes; i++ {
				if msg.String() == fmt.Sprintf("%d", i+1) {
					// First switch focus to this node
					m.selectedNode = i
					// Then toggle the node state
					return m, ToggleNode(i)
				}
			}
		}

	case SwitchFocusMsg:
		m.selectedNode = msg.nodeIdx
		return m, nil

	case ToggleNodeMsg:
		if m.nodeStatuses[msg.nodeIdx] == Running {
			return m, RunNode(msg.nodeIdx, false, m) // Stop node
		}
		return m, RunNode(msg.nodeIdx, true, m) // Start node

	case UpdateStatusMsg:
		m.nodeStatuses[msg.nodeIdx] = msg.status
		return m, UpdateDeemon()

	case UpdateLogsMsg:
		return m, UpdateDeemon()
	}

	return m, nil
}

// View renders the interface.
func (m MultiNode) View() string {
	if m.showHelp {
		return renderHelpView()
	}

	// Create tabs for nodes
	tabs := []string{}
	for i := 0; i < m.numNodes; i++ {
		var status string
		if m.nodeStatuses[i] == Running {
			status = "●"
		} else {
			status = "○"
		}

		tabText := fmt.Sprintf("Node %d %s", i+1, status)

		// apply different styling based on node status and selection
		if i == m.selectedNode {
			if m.nodeStatuses[i] == Running {
				tabs = append(tabs, activeRunningTabStyle.Render(tabText))
			} else {
				tabs = append(tabs, activeTabStyle.Render(tabText))
			}
		} else {
			if m.nodeStatuses[i] == Running {
				tabs = append(tabs, runningTabStyle.Render(tabText))
			} else {
				tabs = append(tabs, tabStyle.Render(tabText))
			}
		}
	}

	// Render the tab row
	tabRow := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)

	// Header row with status
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		headerStyle.Render("Ignite Node Dashboard"),
	)

	// Render selected node details
	nodeDetails := renderNodeDetails(m, m.selectedNode)

	// Render the keyboard controls help at the bottom
	controls := fmt.Sprintf("%s ←/→: Switch node • %s 1-%d: Toggle node • %s q: Quit • %s h: Help",
		infoStyle.Render("•"),
		infoStyle.Render("•"),
		m.numNodes,
		infoStyle.Render("•"),
		infoStyle.Render("•"),
	)

	// Assemble the final view
	return fmt.Sprintf("%s\n%s\n\n%s\n\n%s",
		header,
		tabRow,
		nodeDetails,
		controls,
	)
}

// renderNodeDetails renders the details of a specific node.
func renderNodeDetails(m MultiNode, nodeIdx int) string {
	status := nodeStoppedStyle.Render("[Stopped]")
	statusVerb := "start"

	if m.nodeStatuses[nodeIdx] == Running {
		status = nodeActiveStyle.Render("[Running]")
		statusVerb = "stop"
	}

	tcpAddress := tcpStyle.Render(fmt.Sprintf("tcp://127.0.0.1:%d", m.args.ListPorts[nodeIdx]))
	nodeInfo := fmt.Sprintf("Node %d %s\nEndpoint: %s",
		nodeIdx+1,
		status,
		tcpAddress,
	)

	// Action button
	actionPrompt := fmt.Sprintf("Press [%d] to %s", nodeIdx+1, statusVerb)

	// Log section
	var logContent string
	if len(m.logs[nodeIdx]) > 0 {
		logEntries := []string{}
		for _, line := range m.logs[nodeIdx] {
			logEntries = append(logEntries, logEntryStyle.Render(line))
		}
		logContent = strings.Join(logEntries, "\n")
	} else {
		logContent = infoStyle.Render("No logs available")
	}

	logs := fmt.Sprintf("Logs:\n%s", logBoxStyle.Render(logContent))

	return fmt.Sprintf("%s\n%s\n\n%s", nodeInfo, actionPrompt, logs)
}

// renderHelpView displays help information.
func renderHelpView() string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(subtleColor).
		Padding(1, 2).
		Render(`Ignite Node Dashboard Help

Navigation:
  • Left/Right or Tab/Shift+Tab: Switch between nodes
  • 1-4: Toggle the corresponding node on/off
  • h: Toggle this help screen
  • q or Ctrl+c: Quit and stop all nodes

Node Status:
  • [Running]: The node is active and processing blocks
  • [Stopped]: The node is inactive

This dashboard allows you to manage multiple validator nodes
in your local testnet environment. You can start and stop nodes
independently and monitor their logs in real-time.

Press h to return to the dashboard.`)
}
