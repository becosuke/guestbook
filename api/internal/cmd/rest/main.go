package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pkgconfig "github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/api/internal/registry/injection/rest"
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

	app := rest.InitializeApp(ctx)
	config := app.Config
	logger := app.Logger
	defer func() {
		_ = logger.Sync()
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(interrupt)

	err := pbgo.RegisterGuestbookServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", config.GrpcHost, config.GrpcPort), opts)
	if err != nil {
		logger.Error("rest server: failed to register handler", zap.Error(err))
		return exitError
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.RestHost, config.RestPort),
		Handler: mux,
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("rest server failed to listen", zap.Error(err))
			return err
		}
		return nil
	})
	logger.Info("rest server: listening", zap.String("host", config.RestHost), zap.Int("port", config.RestPort))

	select {
	case <-interrupt:
		logger.Info("received shutdown signal")
	case <-ctx.Done():
		logger.Info("context canceled")
	}
	cancel()

	logger.Info("rest server: going gracefully shutdown")
	if err := httpServer.Shutdown(ctx); err != context.Canceled {
		logger.Error("received error on gracefully shutdown", zap.Error(err))
	}
	logger.Info("rest server: completed gracefully shutdown")

	if err := eg.Wait(); err != nil {
		logger.Error("received error", zap.Error(err))
		return exitError
	}

	return exitOK
}
