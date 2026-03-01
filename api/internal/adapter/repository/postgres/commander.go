package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/domain/repository"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
)

const uniqueViolationCode = "23505"

func NewCommander(config *pkgconfig.Config, logger *zap.Logger, pool *pgxpool.Pool) repository.Commander {
	return &commanderImpl{
		config: config,
		logger: logger,
		pool:   pool,
	}
}

type commanderImpl struct {
	config *pkgconfig.Config
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func (impl *commanderImpl) Create(ctx context.Context, post *domain.Post) error {
	_, err := impl.pool.Exec(ctx,
		"INSERT INTO posts (post_id, body) VALUES ($1, $2)",
		post.PostID().String(), post.Body().String(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return repository.ErrAlreadyExists
		}
		return errors.WithStack(err)
	}
	return nil
}

func (impl *commanderImpl) Update(ctx context.Context, post *domain.Post) error {
	ct, err := impl.pool.Exec(ctx,
		"UPDATE posts SET body = $1, update_time = NOW() WHERE post_id = $2",
		post.Body().String(), post.PostID().String(),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (impl *commanderImpl) Delete(ctx context.Context, postID *domain.PostID) error {
	_, err := impl.pool.Exec(ctx,
		"DELETE FROM posts WHERE post_id = $1",
		postID.String(),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
