package cmdmodel

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/v28/ignite/services/chain"
)

// NodeStatus is an integer data type that represents the status of a node.
type NodeStatus int

const (
	// Stopped indicates that the node is currently stopped.
	Stopped NodeStatus = iota

	// Running indicates that the node is currently running.
	Running
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
	return MultiNode{
		ctx:          ctx,
		appd:         chainname + "d",
		args:         args,
		nodeStatuses: make([]NodeStatus, numNodes), // initial states of nodes
		pids:         make([]int, numNodes),
		numNodes:     numNodes,
		logs:         make([][]string, numNodes), // Initialize logs for each node
	}, nil
}

// Init implements the Init method of the tea.Model interface.
func (m MultiNode) Init() tea.Cmd {
	return nil
}

// ToggleNode toggles the state of a node.
func ToggleNode(nodeIdx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleNodeMsg{nodeIdx: nodeIdx}
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
	for i := 0; i < m.numNodes; i++ {
		if m.nodeStatuses[i] == Running {
			RunNode(i, false, *m)() // Stop node
		}
	}
}

// Update handles messages and updates the model.
func (m MultiNode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.StopAllNodes() // Stop all nodes before quitting
			return m, tea.Quit
		default:
			// Check for numbers from 1 to numNodes
			for i := 0; i < m.numNodes; i++ {
				if msg.String() == fmt.Sprintf("%d", i+1) {
					return m, ToggleNode(i)
				}
			}
		}

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
	// Define styles for the state
	runningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))                               // green
	stoppedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))                               // red
	tcpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))                                   // yellow
	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))                                  // gray
	purpleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))                                // purple
	statusBarStyle := lipgloss.NewStyle().Background(lipgloss.Color("0"))                             // Status bar style
	blueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Background(lipgloss.Color("0")) // blue

	statusBar := blueStyle.Render("Press q to quit | Press 1-4 to ") + statusBarStyle.Render(runningStyle.Render("start")) + blueStyle.Render("/") + statusBarStyle.Render(stoppedStyle.Render("stop")) + blueStyle.Render(" corresponding node")
	output := statusBar + "\n\n"

	// Add node control section
	output += purpleStyle.Render("Node Control:")
	for i := 0; i < m.numNodes; i++ {
		status := stoppedStyle.Render("[Stopped]")
		if m.nodeStatuses[i] == Running {
			status = runningStyle.Render("[Running]")
		}

		tcpAddress := tcpStyle.Render(fmt.Sprintf("tcp://127.0.0.1:%d", m.args.ListPorts[i]))
		nodeGray := grayStyle.Render("--node")
		nodeNumber := purpleStyle.Render(fmt.Sprintf("%d.", i+1))

		output += fmt.Sprintf("\n%s Node %d %s %s %s:\n", nodeNumber, i+1, status, nodeGray, tcpAddress)
		output += " [\n"
		if m.logs != nil {
			for _, line := range m.logs[i] {
				output += "  " + line + "\n"
			}
		}

		output += " ]\n\n"
	}

	output += grayStyle.Render("\nPress q to quit.\n")
	return output
}
