package plugintypes

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type HookMapper interface {
	Hooks() []string
}

type HookMapperRPC struct {
	client *rpc.Client
}

func (c *HookMapperRPC) Hooks() []string {
	var Hooks []string
	cErr := c.client.Call("Plugin.Hooks", new(interface{}), &Hooks)
	if cErr != nil {
		panic(cErr)
	}

	return Hooks
}

type HookMapperRPCServer struct {
	Impl HookMapper
}

func (c *HookMapperRPCServer) Hooks(args interface{}, resp *[]string) error {
	*resp = c.Impl.Hooks()
	return nil
}

type HookMapperPlugin struct {
	Impl HookMapper
}

func (p *HookMapperPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &HookMapperRPCServer{Impl: p.Impl}, nil
}

func (HookMapperPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &HookMapperRPC{client: c}, nil
}

type Hook struct {
	ParentCommand []string
	Name          string
	HookType      string
}

type HookModule interface {
	GetParentCommand() []string
	GetName() string
	GetType() string

	PreRun(*cobra.Command, []string) error
	PostRun(*cobra.Command, []string) error
}

type HookModuleRPC struct {
	client *rpc.Client
}

func (h *HookModuleRPC) GetParentCommand() []string {
	var parentHook []string
	cErr := h.client.Call("Plugin.GetParentCommand", new(interface{}), &parentHook)
	if cErr != nil {
		panic(cErr)
	}

	return parentHook
}

func (h *HookModuleRPC) GetName() string {
	var name string
	cErr := h.client.Call("Plugin.GetName", new(interface{}), &name)
	if cErr != nil {
		panic(cErr)
	}

	return name
}

func (h *HookModuleRPC) GetType() string {
	var hook_type string
	cErr := h.client.Call("Plugin.GetType", new(interface{}), &hook_type)
	if cErr != nil {
		panic(cErr)
	}

	return hook_type
}

func (h *HookModuleRPC) PreRun(cmd *cobra.Command, args []string) error {
	var err error
	cErr := h.client.Call("Plugin.PreRun", ExecArgs{cmd, args}, &err)
	if cErr != nil {
		panic(cErr)
	}

	return err
}

func (h *HookModuleRPC) PostRun(cmd *cobra.Command, args []string) error {
	var err error
	cErr := h.client.Call("Plugin.PostRun", ExecArgs{cmd, args}, &err)
	if cErr != nil {
		panic(cErr)
	}

	return err
}

type HookRPCServer struct {
	Impl HookModule
}

func (h *HookRPCServer) GetParentCommand(args interface{}, resp *[]string) error {
	*resp = h.Impl.GetParentCommand()
	return nil
}

func (h *HookRPCServer) GetName(args interface{}, resp *string) error {
	*resp = h.Impl.GetName()
	return nil
}

func (h *HookRPCServer) GetType(args interface{}, resp *string) error {
	*resp = h.Impl.GetType()
	return nil
}

func (h *HookRPCServer) PreRun(args ExecArgs, resp *error) error {
	*resp = h.Impl.PreRun(args.Cmd, args.Args)
	return nil
}

func (h *HookRPCServer) PostRun(args ExecArgs, resp *error) error {
	*resp = h.Impl.PostRun(args.Cmd, args.Args)
	return nil
}

type HookModulePlugin struct {
	Impl HookModule
}

func (p *HookModulePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &HookRPCServer{Impl: p.Impl}, nil
}

func (HookModulePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &HookModuleRPC{client: c}, nil
}
