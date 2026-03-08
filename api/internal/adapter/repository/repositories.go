package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

type repositoriesImpl struct {
	*postQuerierImpl
	*postCommanderImpl
	*paginatorImpl
}

func NewRepositories(_ context.Context, config *domain.Config, logger *zap.Logger, pool *pgxpool.Pool) interfaces.Repositories {
	return &repositoriesImpl{
		postQuerierImpl:   &postQuerierImpl{config: config, logger: logger, pool: pool},
		postCommanderImpl: &postCommanderImpl{config: config, logger: logger, pool: pool},
		paginatorImpl:     &paginatorImpl{config: config, logger: logger, pool: pool},
	}
}
