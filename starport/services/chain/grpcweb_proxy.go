package chain

import (
	"context"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// newGRPCWebProxyHandler GRPC Web handler to proxy GRPC Web requests to sdk.
// https://github.com/improbable-eng/grpc-web/tree/master/go/grpcwebproxy used as a reference
// to build necessary config.
func newGRPCWebProxyHandler(grpcServerAddress string) (*grpc.ClientConn, http.Handler, error) {
	grpcopts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithCodec(proxy.Codec()),
	}

	grpconn, err := grpc.Dial(grpcServerAddress, grpcopts...)
	if err != nil {
		return nil, nil, err
	}

	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		outCtx, _ := context.WithCancel(ctx)
		mdCopy := md.Copy()
		delete(mdCopy, "user-agent")
		// If this header is present in the request from the web client,
		// the actual connection to the backend will not be established.
		// https://github.com/improbable-eng/grpc-web/issues/568
		delete(mdCopy, "connection")
		outCtx = metadata.NewOutgoingContext(outCtx, mdCopy)
		return outCtx, grpconn, nil
	}

	// Server with logging and monitoring enabled.
	grpcserver := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()), // needed for proxy to function.
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
		grpc_middleware.WithUnaryServerChain(),
		grpc_middleware.WithStreamServerChain(),
	)

	grpcwebopts := []grpcweb.Option{
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool { return true }),
	}

	grpcwebserver := grpcweb.WrapServer(grpcserver, grpcwebopts...)

	return grpconn, grpcwebserver, nil
}
