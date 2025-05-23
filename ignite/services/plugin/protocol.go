package plugin

import (
	"context"
	"sync"

	hplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	v1 "github.com/ignite/cli/v29/ignite/services/plugin/grpc/v1"
)

var handshakeConfig = hplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "IGNITE_APP",
	MagicCookieValue: "ignite",
}

// HandshakeConfig are used to just do a basic handshake between a plugin and host.
// If the handshake fails, a user friendly error is shown. This prevents users from
// executing bad plugins or executing a plugin directory. It is a UX feature, not a
// security feature.
func HandshakeConfig() hplugin.HandshakeConfig {
	return handshakeConfig
}

// NewGRPC returns a new gRPC plugin that implements the interface over gRPC.
func NewGRPC(impl Interface) hplugin.Plugin {
	return grpcPlugin{impl: impl}
}

type grpcPlugin struct {
	hplugin.NetRPCUnsupportedPlugin

	impl Interface
}

// GRPCServer returns a new server that implements the plugin interface over gRPC.
func (p grpcPlugin) GRPCServer(broker *hplugin.GRPCBroker, s *grpc.Server) error {
	v1.RegisterInterfaceServiceServer(s, &server{
		impl:   p.impl,
		broker: broker,
	})
	return nil
}

// GRPCClient returns a new plugin client that allows calling the plugin interface over gRPC.
func (p grpcPlugin) GRPCClient(_ context.Context, broker *hplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &client{
		grpc:   v1.NewInterfaceServiceClient(c),
		broker: broker,
	}, nil
}

type client struct {
	grpc   v1.InterfaceServiceClient
	broker *hplugin.GRPCBroker
}

func (c client) Manifest(ctx context.Context) (*Manifest, error) {
	r, err := c.grpc.Manifest(ctx, &v1.ManifestRequest{})
	if err != nil {
		return nil, err
	}

	return r.Manifest, nil
}

func (c client) Execute(ctx context.Context, cmd *ExecutedCommand, api ClientAPI) error {
	brokerID, stopServer := c.startClientAPIServer(api)
	_, err := c.grpc.Execute(ctx, &v1.ExecuteRequest{
		Cmd:       cmd,
		ClientApi: brokerID,
	})
	stopServer()
	return err
}

func (c client) ExecuteHookPre(ctx context.Context, h *ExecutedHook, api ClientAPI) error {
	brokerID, stopServer := c.startClientAPIServer(api)
	_, err := c.grpc.ExecuteHookPre(ctx, &v1.ExecuteHookPreRequest{
		Hook:      h,
		ClientApi: brokerID,
	})
	stopServer()
	return err
}

func (c client) ExecuteHookPost(ctx context.Context, h *ExecutedHook, api ClientAPI) error {
	brokerID, stopServer := c.startClientAPIServer(api)
	_, err := c.grpc.ExecuteHookPost(ctx, &v1.ExecuteHookPostRequest{
		Hook:      h,
		ClientApi: brokerID,
	})
	stopServer()
	return err
}

func (c client) ExecuteHookCleanUp(ctx context.Context, h *ExecutedHook, api ClientAPI) error {
	brokerID, stopServer := c.startClientAPIServer(api)
	_, err := c.grpc.ExecuteHookCleanUp(ctx, &v1.ExecuteHookCleanUpRequest{
		Hook:      h,
		ClientApi: brokerID,
	})
	stopServer()
	return err
}

func (c client) startClientAPIServer(api ClientAPI) (uint32, func()) {
	var (
		srv      *grpc.Server
		m        sync.Mutex
		brokerID = c.broker.NextId()
	)

	go c.broker.AcceptAndServe(brokerID, func(opts []grpc.ServerOption) *grpc.Server {
		m.Lock()
		defer m.Unlock()
		srv = grpc.NewServer(opts...)
		v1.RegisterClientAPIServiceServer(srv, &clientAPIServer{impl: api})
		return srv
	})

	stop := func() {
		m.Lock()
		if srv != nil {
			srv.Stop()
		}
		m.Unlock()
	}

	return brokerID, stop
}

type server struct {
	v1.UnimplementedInterfaceServiceServer

	impl   Interface
	broker *hplugin.GRPCBroker
}

func (s server) Manifest(ctx context.Context, _ *v1.ManifestRequest) (*v1.ManifestResponse, error) {
	m, err := s.impl.Manifest(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.ManifestResponse{Manifest: m}, nil
}

func (s server) Execute(ctx context.Context, r *v1.ExecuteRequest) (*v1.ExecuteResponse, error) {
	conn, err := s.broker.Dial(r.ClientApi)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = s.impl.Execute(ctx, r.GetCmd(), newClientAPIClient(conn))
	if err != nil {
		return nil, err
	}

	return &v1.ExecuteResponse{}, nil
}

func (s server) ExecuteHookPre(ctx context.Context, r *v1.ExecuteHookPreRequest) (*v1.ExecuteHookPreResponse, error) {
	conn, err := s.broker.Dial(r.ClientApi)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = s.impl.ExecuteHookPre(ctx, r.GetHook(), newClientAPIClient(conn))
	if err != nil {
		return nil, err
	}

	return &v1.ExecuteHookPreResponse{}, nil
}

func (s server) ExecuteHookPost(ctx context.Context, r *v1.ExecuteHookPostRequest) (*v1.ExecuteHookPostResponse, error) {
	conn, err := s.broker.Dial(r.ClientApi)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = s.impl.ExecuteHookPost(ctx, r.GetHook(), newClientAPIClient(conn))
	if err != nil {
		return nil, err
	}

	return &v1.ExecuteHookPostResponse{}, nil
}

func (s server) ExecuteHookCleanUp(ctx context.Context, r *v1.ExecuteHookCleanUpRequest) (*v1.ExecuteHookCleanUpResponse, error) {
	conn, err := s.broker.Dial(r.ClientApi)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = s.impl.ExecuteHookCleanUp(ctx, r.GetHook(), newClientAPIClient(conn))
	if err != nil {
		return nil, err
	}

	return &v1.ExecuteHookCleanUpResponse{}, nil
}

func newClientAPIClient(c *grpc.ClientConn) *clientAPIClient {
	return &clientAPIClient{v1.NewClientAPIServiceClient(c)}
}

type clientAPIClient struct {
	grpc v1.ClientAPIServiceClient
}

func (c clientAPIClient) GetChainInfo(ctx context.Context) (*ChainInfo, error) {
	r, err := c.grpc.GetChainInfo(ctx, &v1.GetChainInfoRequest{})
	if err != nil {
		return nil, err
	}

	return r.ChainInfo, nil
}

func (c clientAPIClient) GetIgniteInfo(ctx context.Context) (*IgniteInfo, error) {
	r, err := c.grpc.GetIgniteInfo(ctx, &v1.GetIgniteInfoRequest{})
	if err != nil {
		return nil, err
	}

	return r.IgniteInfo, nil
}

type clientAPIServer struct {
	v1.UnimplementedClientAPIServiceServer

	impl ClientAPI
}

func (s clientAPIServer) GetChainInfo(ctx context.Context, _ *v1.GetChainInfoRequest) (*v1.GetChainInfoResponse, error) {
	chainInfo, err := s.impl.GetChainInfo(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.GetChainInfoResponse{ChainInfo: chainInfo}, nil
}

func (s clientAPIServer) GetIgniteInfo(ctx context.Context, _ *v1.GetIgniteInfoRequest) (*v1.GetIgniteInfoResponse, error) {
	igniteInfo, err := s.impl.GetIgniteInfo(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.GetIgniteInfoResponse{IgniteInfo: igniteInfo}, nil
}
