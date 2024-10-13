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

	"github.com/ignite/cli/v29/ignite/services/chain"
)

type NodeStatus int

const (
	Stopped NodeStatus = iota
	Running
)

type MultiNode struct {
	ctx  context.Context
	appd string
	args chain.MultiNodeArgs

	nodeStatuses []NodeStatus
	pids         []int      // Store the PIDs of the running processes
	numNodes     int        // Number of nodes
	logs         [][]string // Store logs for each node
}

type ToggleNodeMsg struct {
	nodeIdx int
}

type UpdateStatusMsg struct {
	nodeIdx int
	status  NodeStatus
}

type UpdateLogsMsg struct{}

func UpdateDeemon() tea.Cmd {
	return func() tea.Msg {
		return UpdateLogsMsg{}
	}
}

// NewModel initializes the model.
func NewModel(ctx context.Context, chainname string, args chain.MultiNodeArgs) MultiNode {
	numNodes, err := strconv.Atoi(args.NumValidator)
	if err != nil {
		panic(err)
	}
	return MultiNode{
		ctx:          ctx,
		appd:         chainname + "d",
		args:         args,
		nodeStatuses: make([]NodeStatus, numNodes), // initial states of nodes
		pids:         make([]int, numNodes),
		numNodes:     numNodes,
		logs:         make([][]string, numNodes), // Initialize logs for each node
	}
}

// Implement the Update function.
func (m MultiNode) Init() tea.Cmd {
	return nil
}

// ToggleNode toggles the state of a node.
func ToggleNode(nodeIdx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleNodeMsg{nodeIdx: nodeIdx}
	}
}

// Run or stop the node based on its status.
func RunNode(nodeIdx int, start bool, m MultiNode) tea.Cmd {
	var (
		pid  = &m.pids[nodeIdx]
		args = m.args
		appd = m.appd
	)

	return func() tea.Msg {
		if start {
			nodeHome := filepath.Join(args.OutputDir, args.NodeDirPrefix+strconv.Itoa(nodeIdx))
			// Create the command to run in background as a daemon
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
			go func() {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					line := scanner.Text()
					// Add log line to the respective node's log slice
					m.logs[nodeIdx] = append(m.logs[nodeIdx], line)
					// Keep only the last 5 lines
					if len(m.logs[nodeIdx]) > 5 {
						m.logs[nodeIdx] = m.logs[nodeIdx][len(m.logs[nodeIdx])-5:]
					}
				}
			}()
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

// Stop all nodes.
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
	output := "Node Control:\n"
	for i := 0; i < m.numNodes; i++ {
		status := "[Stopped]"
		if m.nodeStatuses[i] == Running {
			status = "[Running]"
		}
		output += fmt.Sprintf("%d. Node %d %s --node tcp://127.0.0.1:%d:\n", i+1, i+1, status, 26657-3*i)
		output += " [\n"
		for _, line := range m.logs[i] {
			output += "  " + line + "\n"
		}
		output += " ]\n\n"
	}

	output += "Press q to quit.\n"
	return output
}
