package plugin

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	v1 "github.com/ignite/cli/ignite/services/plugin/grpc/v1"
)

var handshakeConfig = hplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
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

func (c client) Execute(ctx context.Context, cmd *ExecutedCommand, a Analizer) error {
	brokerID, stopServer := c.startAnalizerServer(a)
	_, err := c.grpc.Execute(ctx, &v1.ExecuteRequest{
		Cmd:      cmd,
		Analizer: brokerID,
	})
	stopServer()
	return err
}

func (c client) ExecuteHookPre(ctx context.Context, h *ExecutedHook, a Analizer) error {
	brokerID, stopServer := c.startAnalizerServer(a)
	_, err := c.grpc.ExecuteHookPre(ctx, &v1.ExecuteHookPreRequest{
		Hook:     h,
		Analizer: brokerID,
	})
	stopServer()
	return err
}

func (c client) ExecuteHookPost(ctx context.Context, h *ExecutedHook, a Analizer) error {
	brokerID, stopServer := c.startAnalizerServer(a)
	_, err := c.grpc.ExecuteHookPost(ctx, &v1.ExecuteHookPostRequest{
		Hook:     h,
		Analizer: brokerID,
	})
	stopServer()
	return err
}

func (c client) ExecuteHookCleanUp(ctx context.Context, h *ExecutedHook, a Analizer) error {
	brokerID, stopServer := c.startAnalizerServer(a)
	_, err := c.grpc.ExecuteHookCleanUp(ctx, &v1.ExecuteHookCleanUpRequest{
		Hook:     h,
		Analizer: brokerID,
	})
	stopServer()
	return err
}

func (c client) startAnalizerServer(a Analizer) (uint32, func()) {
	var (
		srv      *grpc.Server
		brokerID = c.broker.NextId()
	)

	go c.broker.AcceptAndServe(brokerID, func(opts []grpc.ServerOption) *grpc.Server {
		srv = grpc.NewServer(opts...)
		v1.RegisterAnalizerServiceServer(srv, &analizerServer{impl: a})
		return srv
	})

	return brokerID, func() { srv.Stop() }
}

type server struct {
	v1.UnimplementedInterfaceServiceServer

	impl   Interface
	broker *hplugin.GRPCBroker
}

func (s server) Manifest(ctx context.Context, r *v1.ManifestRequest) (*v1.ManifestResponse, error) {
	m, err := s.impl.Manifest(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.ManifestResponse{Manifest: m}, nil
}

func (s server) Execute(ctx context.Context, r *v1.ExecuteRequest) (*v1.ExecuteResponse, error) {
	conn, err := s.broker.Dial(r.Analizer)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = s.impl.Execute(ctx, r.GetCmd(), newAnalizerClient(conn))
	if err != nil {
		return nil, err
	}

	return &v1.ExecuteResponse{}, nil
}

func (s server) ExecuteHookPre(ctx context.Context, r *v1.ExecuteHookPreRequest) (*v1.ExecuteHookPreResponse, error) {
	conn, err := s.broker.Dial(r.Analizer)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = s.impl.ExecuteHookPre(ctx, r.GetHook(), newAnalizerClient(conn))
	if err != nil {
		return nil, err
	}

	return &v1.ExecuteHookPreResponse{}, nil
}

func (s server) ExecuteHookPost(ctx context.Context, r *v1.ExecuteHookPostRequest) (*v1.ExecuteHookPostResponse, error) {
	conn, err := s.broker.Dial(r.Analizer)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = s.impl.ExecuteHookPost(ctx, r.GetHook(), newAnalizerClient(conn))
	if err != nil {
		return nil, err
	}

	return &v1.ExecuteHookPostResponse{}, nil
}

func (s server) ExecuteHookCleanUp(ctx context.Context, r *v1.ExecuteHookCleanUpRequest) (*v1.ExecuteHookCleanUpResponse, error) {
	conn, err := s.broker.Dial(r.Analizer)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = s.impl.ExecuteHookCleanUp(ctx, r.GetHook(), newAnalizerClient(conn))
	if err != nil {
		return nil, err
	}

	return &v1.ExecuteHookCleanUpResponse{}, nil
}

func newAnalizerClient(c *grpc.ClientConn) *analizerClient {
	return &analizerClient{v1.NewAnalizerServiceClient(c)}
}

type analizerClient struct {
	grpc v1.AnalizerServiceClient
}

func (c analizerClient) Dependencies(ctx context.Context) ([]*Dependency, error) {
	r, err := c.grpc.Dependencies(ctx, &v1.DependenciesRequest{})
	if err != nil {
		return nil, err
	}

	return r.Dependencies, nil
}

type analizerServer struct {
	v1.UnimplementedAnalizerServiceServer

	impl Analizer
}

func (s analizerServer) Dependencies(ctx context.Context, _ *v1.DependenciesRequest) (*v1.DependenciesResponse, error) {
	deps, err := s.impl.Dependencies(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.DependenciesResponse{Dependencies: deps}, nil
}
