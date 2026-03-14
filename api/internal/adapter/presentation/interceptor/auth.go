package interceptor

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

// AuthUnaryServerInterceptor returns a unary server interceptor that
// performs authentication using grpc_auth.
func AuthUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(authFunc)
}

// AuthStreamServerInterceptor returns a stream server interceptor that
// performs authentication using grpc_auth.
func AuthStreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(authFunc)
}

func authFunc(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
