package pluginsrpc

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type Command interface {
	ParentCommand() []string
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
	NumArgs() int
	Exec(*cobra.Command, []string) error
}

type CommandModule interface {
	Module
	Registry() map[string]Command
}

type CommandModuleRPC struct {
	client *rpc.Client
}

func (c *CommandModuleRPC) Init(ctx context.Context) error {
	var err error
	cErr := c.client.Call("Plugin.ParentCommand", InitArgs{ctx}, &err)
	if cErr != nil {
		panic(cErr)
	}

	return err
}

func (c *CommandModuleRPC) Registry(ctx context.Context) map[string]Command {
	var commands map[string]Command
	cErr := c.client.Call("Plugin.Registry", new(interface{}), &commands)
	if cErr != nil {
		panic(cErr)
	}

	return commands
}

// func (h *CommandRPC) ParentCommand() []string {
// 	var parentCommand []string
// 	err := h.client.Call("Plugin.ParentCommand", new(interface{}), &parentCommand)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return parentCommand
// }

// func (h *CommandRPC) Name() string {
// 	var name string
// 	err := h.client.Call("Plugin.Name", new(interface{}), &name)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return name
// }

// func (h *CommandRPC) Usage() string {
// 	var commandUsage string
// 	err := h.client.Call("Plugin.Usage", new(interface{}), &commandUsage)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return commandUsage
// }

// func (h *CommandRPC) ShortDesc() string {
// 	var shortDesc string
// 	err := h.client.Call("Plugin.ShortDesc", new(interface{}), &shortDesc)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return shortDesc
// }

// func (h *CommandRPC) LongDesc() string {
// 	var longDesc string
// 	err := h.client.Call("Plugin.LongDesc", new(interface{}), &longDesc)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return longDesc
// }

// func (h *CommandRPC) NumArgs() int {
// 	var numArgs int
// 	cErr := h.client.Call("Plugin.NumArgs", new(interface{}), &numArgs)
// 	if cErr != nil {
// 		panic(cErr)
// 	}

// 	return numArgs
// }

// func (h *CommandRPC) Exec(cmd *cobra.Command, args []string) error {
// 	var err error
// 	cErr := h.client.Call("Plugin.PreRun", PostRunArgs{
// 		cmd:  cmd,
// 		args: args,
// 	}, &err)
// 	if cErr != nil {
// 		panic(cErr)
// 	}

// 	return err
// }

type CommandRPCServer struct {
	Impl CommandModule
}

func (h *CommandRPCServer) Init(args InitArgs, resp *error) error {
	*resp = h.Impl.Init(args.ctx)
	return nil
}

func (h *CommandRPCServer) Registry(args interface{}, resp *map[string]Command) error {
	*resp = h.Impl.Registry()
	return nil
}

// func (h *CommandRPCServer) ParentCommand(args interface{}, resp *[]string) error {
// 	*resp = h.Impl.ParentCommand()
// 	return nil
// }

// func (h *CommandRPCServer) Name(args interface{}, resp *string) error {
// 	*resp = h.Impl.Name()
// 	return nil
// }

// func (h *CommandRPCServer) Usage(args interface{}, resp *string) error {
// 	*resp = h.Impl.Usage()
// 	return nil
// }

// func (h *CommandRPCServer) ShortDesc(args interface{}, resp *string) error {
// 	*resp = h.Impl.ShortDesc()
// 	return nil
// }

// func (h *CommandRPCServer) LongDesc(args interface{}, resp *string) error {
// 	*resp = h.Impl.LongDesc()
// 	return nil
// }

// func (h *CommandRPCServer) NumArgs(args interface{}, resp *int) error {
// 	*resp = h.Impl.NumArgs()
// 	return nil
// }

// type ExecArgs struct {
// 	cmd  *cobra.Command
// 	args []string
// }

// func (h *CommandRPCServer) Exec(args ExecArgs, resp *error) error {
// 	*resp = h.Impl.Exec(args.cmd, args.args)
// 	return nil
// }

type CommandPlugin struct {
	Impl CommandModule
}

func (p *CommandPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CommandRPCServer{Impl: p.Impl}, nil
}

func (CommandPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CommandModuleRPC{client: c}, nil
}
