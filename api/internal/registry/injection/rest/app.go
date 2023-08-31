package rest

import (
	"go.uber.org/zap"

	pkgconfig "github.com/becosuke/guestbook/api/internal/registry/config"
)

type App struct {
	Config *pkgconfig.Config
	Logger *zap.Logger
}
