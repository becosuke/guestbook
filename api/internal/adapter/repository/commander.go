package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

const uniqueViolationCode = "23505"

func NewCommander(config *domain.Config, logger *zap.Logger, pool *pgxpool.Pool) interfaces.Commander {
	return &commanderImpl{
		config: config,
		logger: logger,
		pool:   pool,
	}
}

type commanderImpl struct {
	config *domain.Config
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func (impl *commanderImpl) Create(ctx context.Context, post *domain.Post) error {
	_, err := impl.pool.Exec(ctx,
		`INSERT INTO Posts (PostId, PostBody) VALUES ($1, $2)`,
		post.PostID().String(), post.PostBody().String(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (impl *commanderImpl) Update(ctx context.Context, post *domain.Post) error {
	ct, err := impl.pool.Exec(ctx,
		`UPDATE Posts SET PostBody = $1, UpdateTime = NOW() WHERE PostId = $2 AND DeleteTime IS NULL`,
		post.PostBody().String(), post.PostID().String(),
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (impl *commanderImpl) Delete(ctx context.Context, postID *domain.PostID) error {
	ct, err := impl.pool.Exec(ctx,
		`UPDATE Posts SET DeleteTime = NOW(), UpdateTime = NOW() WHERE PostId = $1 AND DeleteTime IS NULL`,
		postID.String(),
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
