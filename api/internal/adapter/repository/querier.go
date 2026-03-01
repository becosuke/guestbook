package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

func NewQuerier(config *domain.Config, logger *zap.Logger, pool *pgxpool.Pool) interfaces.Querier {
	return &querierImpl{
		config: config,
		logger: logger,
		pool:   pool,
	}
}

type querierImpl struct {
	config *domain.Config
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func (impl *querierImpl) Get(ctx context.Context, postID *domain.PostID) (*domain.Post, error) {
	var body string
	err := impl.pool.QueryRow(ctx, `SELECT PostBody FROM Posts WHERE PostId = $1`, postID.String()).Scan(&body)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return domain.NewPost(postID, domain.NewPostBody(body)), nil
}

func (impl *querierImpl) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, error) {
	pageSize := int32(10)
	if pageOption.PageSize() != nil {
		pageSize = int32(*pageOption.PageSize())
	}

	var rows pgx.Rows
	var err error
	if pageOption.PageToken() != nil && string(*pageOption.PageToken()) != "" {
		rows, err = impl.pool.Query(ctx,
			`SELECT PostId, PostBody FROM Posts WHERE PostId < $1 ORDER BY PostId DESC LIMIT $2`,
			string(*pageOption.PageToken()), pageSize,
		)
	} else {
		rows, err = impl.pool.Query(ctx,
			`SELECT PostId, PostBody FROM Posts ORDER BY PostId DESC LIMIT $1`,
			pageSize,
		)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var postID, body string
		if err := rows.Scan(&postID, &body); err != nil {
			return nil, errors.WithStack(err)
		}
		posts = append(posts, domain.NewPost(domain.NewPostID(postID), domain.NewPostBody(body)))
	}
	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	if posts == nil {
		posts = []*domain.Post{}
	}
	return posts, nil
}
