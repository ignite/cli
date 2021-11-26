package pluginsrpc

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type Hook interface {
	ParentCommand() []string
	Name() string
	Type() string
	ShortDesc() string

	PreRun(*cobra.Command, []string) error
	PostRun(*cobra.Command, []string) error
}

type HookModule interface {
	Module
	Registry() map[string]Hook
}

type HookModuleRPC struct {
	client *rpc.Client
}

func (h *HookModuleRPC) Init(ctx context.Context) error {
	var err error
	cErr := h.client.Call("Plugin.ParentCommand", InitArgs{ctx}, &err)
	if cErr != nil {
		panic(cErr)
	}

	return err
}

func (h *HookModuleRPC) Registry(ctx context.Context) map[string]Command {
	var commands map[string]Command
	cErr := h.client.Call("Plugin.Registry", new(interface{}), &commands)
	if cErr != nil {
		panic(cErr)
	}

	return commands
}

// func (h *HookRPC) ParentCommand() []string {
// 	var parentCommand []string
// 	err := h.client.Call("Plugin.ParentCommand", new(interface{}), &parentCommand)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return parentCommand
// }

// func (h *HookRPC) Name() string {
// 	var name string
// 	err := h.client.Call("Plugin.Name", new(interface{}), &name)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return name
// }

// func (h *HookRPC) Type() string {
// 	var hookType string
// 	err := h.client.Call("Plugin.Type", new(interface{}), &hookType)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return hookType
// }

// func (h *HookRPC) ShortDesc() string {
// 	var shortDesc string
// 	err := h.client.Call("Plugin.ParentCommand", new(interface{}), &shortDesc)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return shortDesc
// }

// func (h *HookRPC) PreRun(cmd *cobra.Command, args []string) error {
// 	var err error
// 	cErr := h.client.Call("Plugin.PreRun", PreRunArgs{
// 		cmd:  cmd,
// 		args: args,
// 	}, &err)
// 	if cErr != nil {
// 		panic(cErr)
// 	}

// 	return err
// }

// func (h *HookRPC) PostRun(cmd *cobra.Command, args []string) error {
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

type HookRPCServer struct {
	Impl HookModule
}

func (h *HookRPCServer) Init(args InitArgs, resp *error) error {
	*resp = h.Impl.Init(args.ctx)
	return nil
}

func (h *HookRPCServer) Registry(args interface{}, resp *map[string]Hook) error {
	*resp = h.Impl.Registry()
	return nil
}

// func (h *HookRPCServer) ParentCommand(args interface{}, resp *[]string) error {
// 	*resp = h.Impl.ParentCommand()
// 	return nil
// }

// func (h *HookRPCServer) Name(args interface{}, resp *string) error {
// 	*resp = h.Impl.Name()
// 	return nil
// }

// func (h *HookRPCServer) Type(args interface{}, resp *string) error {
// 	*resp = h.Impl.Type()
// 	return nil
// }

// func (h *HookRPCServer) ShortDesc(args interface{}, resp *string) error {
// 	*resp = h.Impl.ShortDesc()
// 	return nil
// }

// type PreRunArgs struct {
// 	cmd  *cobra.Command
// 	args []string
// }

// func (h *HookRPCServer) PreRun(args PreRunArgs, resp *error) error {
// 	*resp = h.Impl.PreRun(args.cmd, args.args)
// 	return nil
// }

// type PostRunArgs struct {
// 	cmd  *cobra.Command
// 	args []string
// }

// func (h *HookRPCServer) PostRun(args PostRunArgs, resp *error) error {
// 	*resp = h.Impl.PostRun(args.cmd, args.args)
// 	return nil
// }

type HookPlugin struct {
	Impl HookModule
}

func (p *HookPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &HookRPCServer{Impl: p.Impl}, nil
}

func (HookPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &HookModuleRPC{client: c}, nil
}
