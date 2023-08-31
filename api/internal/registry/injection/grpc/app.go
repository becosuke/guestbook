package grpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pkgconfig "github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/pbgo"
)

type App struct {
	Config     *pkgconfig.Config
	Logger     *zap.Logger
	GrpcServer *grpc.Server
	Controller pbgo.GuestbookServiceServer
}
