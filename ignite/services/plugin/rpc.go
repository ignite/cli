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

// NewRPCPlugin returns a new plugin that implements the interface over RPC.
func NewRPCPlugin(impl Interface) plugin.Plugin {
	return rpcPlugin{impl}
}

type rpcPlugin struct {
	impl Interface
}

// Server returns a new RPC server that implements the plugin interface.
func (p rpcPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return rpcServer{p.impl}, nil
}

// Client returns a new RPC client that allows calling the plugin interface over RPC.
func (rpcPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return rpcClient{client: c}, nil
}

type rpcClient struct{ client *rpc.Client }

func (g rpcClient) Manifest() (Manifest, error) {
	var resp Manifest
	return resp, g.client.Call("Plugin.Manifest", new(interface{}), &resp)
}

func (g rpcClient) Execute(c ExecutedCommand) error {
	var resp interface{}
	return g.client.Call("Plugin.Execute", map[string]interface{}{
		"executedCommand": c,
	}, &resp)
}

func (g rpcClient) ExecuteHookPre(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookPre", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

func (g rpcClient) ExecuteHookPost(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookPost", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

func (g rpcClient) ExecuteHookCleanUp(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookCleanUp", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

type rpcServer struct {
	impl Interface
}

func (s rpcServer) Manifest(_ interface{}, resp *Manifest) (err error) {
	*resp, err = s.impl.Manifest()
	return err
}

func (s rpcServer) Execute(args map[string]interface{}, _ *interface{}) error {
	return s.impl.Execute(args["executedCommand"].(ExecutedCommand))
}

func (s rpcServer) ExecuteHookPre(args map[string]interface{}, _ *interface{}) error {
	return s.impl.ExecuteHookPre(args["executedHook"].(ExecutedHook))
}

func (s rpcServer) ExecuteHookPost(args map[string]interface{}, _ *interface{}) error {
	return s.impl.ExecuteHookPost(args["executedHook"].(ExecutedHook))
}

func (s rpcServer) ExecuteHookCleanUp(args map[string]interface{}, _ *interface{}) error {
	return s.impl.ExecuteHookCleanUp(args["executedHook"].(ExecutedHook))
}
