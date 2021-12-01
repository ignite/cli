package plugintypes

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type CommandMapper interface {
	Commands() []string
}

type CommandMapperRPC struct {
	client *rpc.Client
}

func (c *CommandMapperRPC) Commands() []string {
	var commands []string
	cErr := c.client.Call("Plugin.Commands", new(interface{}), &commands)
	if cErr != nil {
		panic(cErr)
	}

	return commands
}

type CommandMapperRPCServer struct {
	Impl CommandMapper
}

func (c *CommandMapperRPCServer) Commands(args interface{}, resp *[]string) error {
	*resp = c.Impl.Commands()
	return nil
}

type CommandMapperPlugin struct {
	Impl CommandMapper
}

func (p *CommandMapperPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CommandMapperRPCServer{Impl: p.Impl}, nil
}

func (CommandMapperPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CommandMapperRPC{client: c}, nil
}

type Command struct {
	ParentCommand []string
	Name          string
	Usage         string
	ShortDesc     string
	LongDesc      string
	NumArgs       int
}

// need getter method for the parameters!
type CommandModule interface {
	GetParentCommand() []string
	GetName() string
	GetUsage() string
	GetShortDesc() string
	GetLongDesc() string
	GetNumArgs() int
	Exec(*cobra.Command, []string) error
}

type CommandModuleRPC struct {
	client *rpc.Client
}

func (c *CommandModuleRPC) GetParentCommand() []string {
	var commands []string
	cErr := c.client.Call("Plugin.GetParentCommand", new(interface{}), &commands)
	if cErr != nil {
		panic(cErr)
	}

	return commands
}

func (c *CommandModuleRPC) GetName() string {
	var name string
	cErr := c.client.Call("Plugin.GetName", new(interface{}), &name)
	if cErr != nil {
		panic(cErr)
	}

	return name
}

func (c *CommandModuleRPC) GetUsage() string {
	var usage string
	cErr := c.client.Call("Plugin.GetUsage", new(interface{}), &usage)
	if cErr != nil {
		panic(cErr)
	}

	return usage
}

func (c *CommandModuleRPC) GetShortDesc() string {
	var desc string
	cErr := c.client.Call("Plugin.GetShortDesc", new(interface{}), &desc)
	if cErr != nil {
		panic(cErr)
	}

	return desc
}

func (c *CommandModuleRPC) GetLongDesc() string {
	var desc string
	cErr := c.client.Call("Plugin.GetLongDesc", new(interface{}), &desc)
	if cErr != nil {
		panic(cErr)
	}

	return desc
}

func (c *CommandModuleRPC) GetNumArgs() int {
	var numArgs int
	cErr := c.client.Call("Plugin.GetNumArgs", new(interface{}), &numArgs)
	if cErr != nil {
		panic(cErr)
	}

	return numArgs
}

func (c *CommandModuleRPC) Exec(cmd *cobra.Command, args []string) error {
	var err error
	cErr := c.client.Call("Plugin.Exec", ExecArgs{cmd, args}, &err)
	if cErr != nil {
		panic(cErr)
	}

	return err
}

type CommandModuleRPCServer struct {
	Impl CommandModule
}

func (c *CommandModuleRPCServer) GetParentCommand(args interface{}, resp *[]string) error {
	*resp = c.Impl.GetParentCommand()
	return nil
}

func (c *CommandModuleRPCServer) GetName(args interface{}, resp *string) error {
	*resp = c.Impl.GetName()
	return nil
}

func (c *CommandModuleRPCServer) GetUsage(args interface{}, resp *string) error {
	*resp = c.Impl.GetUsage()
	return nil
}

func (c *CommandModuleRPCServer) GetShortDesc(args interface{}, resp *string) error {
	*resp = c.Impl.GetShortDesc()
	return nil
}

func (c *CommandModuleRPCServer) GetLongDesc(args interface{}, resp *string) error {
	*resp = c.Impl.GetLongDesc()
	return nil
}

func (c *CommandModuleRPCServer) GetNumArgs(args interface{}, resp *int) error {
	*resp = c.Impl.GetNumArgs()
	return nil
}

func (c *CommandModuleRPCServer) Exec(args ExecArgs, resp *error) error {
	*resp = c.Impl.Exec(args.Cmd, args.Args)
	return nil
}

type CommandModulePlugin struct {
	Impl CommandModule
}

func (p *CommandModulePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CommandModuleRPCServer{Impl: p.Impl}, nil
}

func (CommandModulePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CommandModuleRPC{client: c}, nil
}
