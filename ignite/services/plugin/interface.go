package plugin

import (
	"encoding/gob"
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

func init() {
	gob.Register(Command{})
}

// An ignite plugin must implements the Plugin interface.
type Interface interface {
	// Commands returns one or multiple commands that will be added to the list
	// of ignite commands. It's invoked each time ignite is executed, in
	// order to display the list of available commands.
	// Each commands are independent, for nested commands, use the field
	// Command.Commands.
	Commands() []Command
	// Execute will be invoked by ignite when a plugin commands is executed.
	// cmd is the executed command (one of the those returned by Commands method)
	// args is the command line arguments passed behing the command.
	Execute(cmd Command, args []string) error
}

// Command represents a plugin command.
type Command struct {
	// Same as cobra.Command.Use
	Use string
	// Same as cobra.Command.Short
	Short string
	// Same as cobra.Command.Long
	Long  string
	Flags []Flag
	// PlaceCommandUnder indicates where the command should be placed.
	// For instance `ignite scaffold` will place the command at the
	// `scaffold` command.
	// An empty value is interpreted as `ignite` (==root).
	PlaceCommandUnder string
	// List of sub commands
	Commands []Command

	// The following fields are populated at runtime
	CobraCmd *cobra.Command
	// Optionnal parameters populated by config at runtime via
	// chainconfig.Plugin.With field.
	With map[string]string
}

type Flag struct {
	Name      string // name as it appears on command line
	Shorthand string // one-letter abbreviated flag
	Usage     string // help message
	DefValue  string // default value (as text); for usage message
	Type      FlagType
}

type FlagType int

const (
	FlagTypeString FlagType = iota
	FlagTypeInt
	FlagTypeBool
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func HandshakeConfig() plugin.HandshakeConfig {
	return handshakeConfig
}

// Here is an implementation that talks over RPC
type InterfaceRPC struct{ client *rpc.Client }

// Commands implements Interface.Commands
func (g *InterfaceRPC) Commands() []Command {
	var resp []Command
	err := g.client.Call("Plugin.Commands", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		log.Fatalf("error while calling plugin %v", err)
	}
	return resp
}

// Execute implements Interface.Commands
func (g *InterfaceRPC) Execute(c Command, args []string) error {
	var resp interface{}
	return g.client.Call("Plugin.Execute", map[string]interface{}{
		"command": c,
		"args":    args,
	}, &resp)
}

// Here is the RPC server that InterfaceRPC talks to, conforming to
// the requirements of net/rpc
type InterfaceRPCServer struct {
	// This is the real implementation
	Impl Interface
}

func (s *InterfaceRPCServer) Commands(args interface{}, resp *[]Command) error {
	*resp = s.Impl.Commands()
	return nil
}

func (s *InterfaceRPCServer) Execute(args map[string]interface{}, resp *interface{}) error {
	return s.Impl.Execute(args["command"].(Command), args["args"].([]string))
}

// This is the implementation of plugin.Interface so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a InterfaceRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return InterfaceRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type InterfacePlugin struct {
	// Impl Injection
	Impl Interface
}

func (p *InterfacePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &InterfaceRPCServer{Impl: p.Impl}, nil
}

func (InterfacePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &InterfaceRPC{client: c}, nil
}
