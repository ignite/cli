package cmdmodel

import (
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
	appd string
	args chain.MultiNodeArgs
	ctx  context.Context

	nodeStatuses []NodeStatus
	pids         []int // Store the PIDs of the running processes
	numNodes     int   // Number of nodes
}

type ToggleNodeMsg struct {
	nodeIdx int
}

type UpdateStatusMsg struct {
	nodeIdx int
	status  NodeStatus
}

// Initialize the model
func NewModel(chainname string, ctx context.Context, args chain.MultiNodeArgs) MultiNode {
	numNodes, err := strconv.Atoi(args.NumValidator)
	if err != nil {
		panic(err)
	}
	return MultiNode{
		appd:         chainname + "d",
		args:         args,
		ctx:          ctx,
		nodeStatuses: make([]NodeStatus, numNodes), // initial states of nodes
		pids:         make([]int, numNodes),
		numNodes:     numNodes,
	}
}

// Implement the Update function
func (m MultiNode) Init() tea.Cmd {
	return nil
}

// ToggleNode toggles the state of a node
func ToggleNode(nodeIdx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleNodeMsg{nodeIdx: nodeIdx}
	}
}

// Run or stop the node based on its status
func RunNode(nodeIdx int, start bool, pid *int, args chain.MultiNodeArgs, appd string) tea.Cmd {
	return func() tea.Msg {
		if start {
			nodeHome := filepath.Join(args.OutputDir, args.NodeDirPrefix+strconv.Itoa(nodeIdx))
			// Create the command to run in background as a daemon
			cmd := exec.Command(appd, "start", "--home", nodeHome)

			// Start the process as a daemon
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setpgid: true, // Ensure it runs in a new process group
			}

			err := cmd.Start() // Start the node in the background
			if err != nil {
				fmt.Printf("Failed to start node %d: %v\n", nodeIdx+1, err)
				return UpdateStatusMsg{nodeIdx: nodeIdx, status: Stopped}
			}

			*pid = cmd.Process.Pid // Store the PID
			go cmd.Wait()          // Let the process run asynchronously
			return UpdateStatusMsg{nodeIdx: nodeIdx, status: Running}
		} else {
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
}

// Stop all nodes
func (m *MultiNode) StopAllNodes() {
	for i := 0; i < m.numNodes; i++ {
		if m.nodeStatuses[i] == Running {
			RunNode(i, false, &m.pids[i], m.args, m.appd)() // Stop node
		}
	}
}

// Update handles messages and updates the model
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
			return m, RunNode(msg.nodeIdx, false, &m.pids[msg.nodeIdx], m.args, m.appd) // Stop node
		}
		return m, RunNode(msg.nodeIdx, true, &m.pids[msg.nodeIdx], m.args, m.appd) // Start node

	case UpdateStatusMsg:
		m.nodeStatuses[msg.nodeIdx] = msg.status
		return m, nil
	}

	return m, nil
}

// View renders the interface
func (m MultiNode) View() string {
	statusText := func(status NodeStatus) string {
		if status == Running {
			return "[Running]"
		}
		return "[Stopped]"
	}

	infoNode := func(i int) string {
		chainId := m.args.ChainID
		home := m.args.OutputDir
		ipaddr := "tcp://127.0.0.1:" + strconv.Itoa(26657-3*i)
		return fmt.Sprintf("INFO: ChainID:%s | Home:%s | Node:%s ", chainId, home, ipaddr)
	}

	output := "Press keys 1,2,3.. to start and stop node 1,2,3.. respectively \nNode Control:\n"
	for i := 0; i < m.numNodes; i++ {
		output += fmt.Sprintf("%d. Node %d %s -- %s\n", i+1, i+1, statusText(m.nodeStatuses[i]), infoNode(i))
	}
	output += "Press q to quit.\n"
	return output
}
