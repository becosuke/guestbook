package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

func NewPaginator(config *domain.Config, logger *zap.Logger, pool *pgxpool.Pool) interfaces.Paginator {
	return &paginatorImpl{
		config: config,
		logger: logger,
		pool:   pool,
	}
}

type paginatorImpl struct {
	config *domain.Config
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func (impl *paginatorImpl) GetPagination(ctx context.Context, paginationID domain.PaginationID) (*domain.Pagination, error) {
	var cursor []byte
	err := impl.pool.QueryRow(ctx, `SELECT Cursor FROM Paginations WHERE PaginationId = $1`, paginationID.String()).Scan(&cursor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return domain.NewPagination(paginationID, cursor), nil
}

func (impl *paginatorImpl) SavePagination(ctx context.Context, pagination *domain.Pagination) error {
	_, err := impl.pool.Exec(ctx,
		`INSERT INTO Paginations (PaginationId, Cursor) VALUES ($1, $2)`,
		pagination.PaginationID().String(), pagination.Cursor(),
	)
	return err
}
