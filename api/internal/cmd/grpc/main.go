package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/reflection"

	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/pbgo"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, pkgconfig.ServiceName{}, serviceName)
	ctx = context.WithValue(ctx, pkgconfig.ServiceVersion{}, version)

	app := InitializeApp(ctx)
	config := app.Config
	logger := app.Logger
	defer func() {
		_ = logger.Sync()
	}()

	grpcServer := app.GrpcServer
	controller := app.Controller
	pbgo.RegisterGuestbookServiceServer(grpcServer, controller)
	reflection.Register(grpcServer)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(interrupt)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GrpcHost, config.GrpcPort))
	if err != nil {
		logger.Error("grpc server: failed to listen", zap.Error(err))
		return exitError
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return grpcServer.Serve(listener)
	})
	logger.Info("grpc server: serving", zap.String("host", config.GrpcHost), zap.Int("port", config.GrpcPort))

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
