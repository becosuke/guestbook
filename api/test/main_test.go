package test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"buf.build/go/protovalidate"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/becosuke/guestbook/api/internal/adapter/presentation/interceptor"
	"github.com/becosuke/guestbook/api/internal/adapter/presentation"
	"github.com/becosuke/guestbook/api/internal/adapter/repository"
	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
	"github.com/becosuke/guestbook/api/internal/usecase"
)

const bufSize = 1024 * 1024

var (
	testClient pb.GuestbookServiceClient
	testPool   *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://guestbook:guestbook@localhost:5432/guestbook?sslmode=disable"
	}

	pool, err := repository.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	testPool = pool

	cfg := &domain.Config{
		Environment: domain.EnvTest,
		LogLevel:    zap.InfoLevel,
		GrpcHost:    "",
		GrpcPort:    50051,
		DatabaseURL: databaseURL,
	}

	zapLogger, _ := zap.NewDevelopment()

	repos := repository.NewRepositories(ctx, cfg, zapLogger, pool)
	uc := usecase.NewUsecase(cfg, zapLogger, repos)
	ctrl := presentation.NewGuestbookServiceServer(cfg, zapLogger, uc)

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to create protovalidate validator: %v", err)
	}

	lis := bufconn.Listen(bufSize)
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zapLogger),
			interceptor.ProtovalidateStreamServerInterceptor(validator),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zapLogger),
			interceptor.ProtovalidateUnaryServerInterceptor(validator),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	pb.RegisterGuestbookServiceServer(server, ctrl)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("server exited with error: %v", err)
		}
	}()

	conn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}

	testClient = pb.NewGuestbookServiceClient(conn)

	code := m.Run()

	conn.Close()
	server.GracefulStop()
	pool.Close()

	os.Exit(code)
}
