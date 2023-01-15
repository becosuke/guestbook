package main

import (
	"context"
	"fmt"
	"github.com/becosuke/guestbook/api/internal/registry/injection"
	"github.com/becosuke/guestbook/api/pb"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	exitOK int = iota
	exitError
)

var (
	serviceName string
	version     string
)

func main() {
	os.Exit(run())
}

func run() int {
	in := injection.NewInjection(serviceName, version)
	config := in.InjectConfig()
	logger := in.InjectLogger()
	defer logger.Sync()

	grpcServer := in.InjectGrpcServer()
	controller := in.InjectController()
	pb.RegisterGuestbookServiceServer(grpcServer, controller)
	reflection.Register(grpcServer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(interrupt)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	if err != nil {
		logger.Error("grpc server: failed to listen", zap.Error(err))
		return exitError
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return grpcServer.Serve(listener)
	})
	logger.Info("grpc server: serving", zap.Int("port", config.GrpcPort))

	select {
	case <-interrupt:
		logger.Info("received shutdown signal")
	case <-ctx.Done():
		logger.Info("context canceled")
	}
	cancel()

	logger.Info("grpc server: going gracefully shutdown")
	grpcServer.GracefulStop()
	logger.Info("grpc server: completed gracefully shutdown")

	if err := eg.Wait(); err != nil {
		logger.Error("received error", zap.Error(err))
		return exitError
	}

	return exitOK
}
