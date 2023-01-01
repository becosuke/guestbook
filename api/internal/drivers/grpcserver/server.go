package grpcserver

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func NewGrpcServer() *grpc.Server {
	return grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(),
		),
	)
}
