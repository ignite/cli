package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// HandshakeConfig are used to just do a basic handshake between a plugin and host.
// If the handshake fails, a user friendly error is shown. This prevents users from
// executing bad plugins or executing a plugin directory. It is a UX feature, not a
// security feature.
func HandshakeConfig() plugin.HandshakeConfig {
	return handshakeConfig
}

// InterfaceRPC is an implementation that talks over RPC.
type InterfaceRPC struct{ client *rpc.Client }

// Manifest implements Interface.Manifest.
func (g *InterfaceRPC) Manifest() (Manifest, error) {
	var resp Manifest
	return resp, g.client.Call("Plugin.Manifest", new(interface{}), &resp)
}

// Execute implements Interface.Commands.
func (g *InterfaceRPC) Execute(c ExecutedCommand) error {
	var resp interface{}
	return g.client.Call("Plugin.Execute", map[string]interface{}{
		"executedCommand": c,
	}, &resp)
}

func (g *InterfaceRPC) ExecuteHookPre(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookPre", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

func (g *InterfaceRPC) ExecuteHookPost(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookPost", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

func (g *InterfaceRPC) ExecuteHookCleanUp(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookCleanUp", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

// InterfaceRPCServer is the RPC server that InterfaceRPC talks to, conforming to
// the requirements of net/rpc.
type InterfaceRPCServer struct {
	// This is the real implementation
	Impl Interface
}

func (s *InterfaceRPCServer) Manifest(_ interface{}, resp *Manifest) error {
	var err error
	*resp, err = s.Impl.Manifest()
	return err
}

func (s *InterfaceRPCServer) Execute(args map[string]interface{}, _ *interface{}) error {
	return s.Impl.Execute(args["executedCommand"].(ExecutedCommand))
}

func (s *InterfaceRPCServer) ExecuteHookPre(args map[string]interface{}, _ *interface{}) error {
	return s.Impl.ExecuteHookPre(args["executedHook"].(ExecutedHook))
}

func (s *InterfaceRPCServer) ExecuteHookPost(args map[string]interface{}, _ *interface{}) error {
	return s.Impl.ExecuteHookPost(args["executedHook"].(ExecutedHook))
}

func (s *InterfaceRPCServer) ExecuteHookCleanUp(args map[string]interface{}, _ *interface{}) error {
	return s.Impl.ExecuteHookCleanUp(args["executedHook"].(ExecutedHook))
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

func (InterfacePlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &InterfaceRPC{client: c}, nil
}
