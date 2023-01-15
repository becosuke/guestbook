package main

import (
	"context"
	"fmt"
	"github.com/becosuke/guestbook/api/internal/registry/injection"
	"github.com/becosuke/guestbook/api/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
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

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(interrupt)

	err := pb.RegisterGuestbookServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.GrpcPort), opts)
	if err != nil {
		logger.Error("http server: failed to register handler", zap.Error(err))
		return exitError
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: mux,
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("http server failed to listen", zap.Error(err))
			return err
		}
		return nil
	})
	logger.Info("http server: listening", zap.Int("port", config.HttpPort))

	select {
	case <-interrupt:
		logger.Info("received shutdown signal")
	case <-ctx.Done():
		logger.Info("context canceled")
	}
	cancel()

	logger.Info("http server: going gracefully shutdown")
	if err := httpServer.Shutdown(ctx); err != context.Canceled {
		logger.Error("received error on gracefully shutdown", zap.Error(err))
	}
	logger.Info("http server: completed gracefully shutdown")

	if err := eg.Wait(); err != nil {
		logger.Error("received error", zap.Error(err))
		return exitError
	}

	return exitOK
}
