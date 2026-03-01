package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain/entity"
	domainrepo "github.com/becosuke/guestbook/api/internal/domain/repository"
)

func NewQuerier(config *entity.Config, logger *zap.Logger, pool *pgxpool.Pool) domainrepo.Querier {
	return &querierImpl{
		config: config,
		logger: logger,
		pool:   pool,
	}
}

type querierImpl struct {
	config *entity.Config
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func (impl *querierImpl) Get(ctx context.Context, postID *entity.PostID) (*entity.Post, error) {
	var body string
	err := impl.pool.QueryRow(ctx, "SELECT body FROM posts WHERE post_id = $1", postID.String()).Scan(&body)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainrepo.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return entity.NewPost(postID, entity.NewPostBody(body)), nil
}

func (impl *querierImpl) Range(ctx context.Context, pageOption *entity.PageOption) ([]*entity.Post, error) {
	pageSize := int32(10)
	if pageOption.PageSize() != nil {
		pageSize = int32(*pageOption.PageSize())
	}

	var rows pgx.Rows
	var err error
	if pageOption.PageToken() != nil && string(*pageOption.PageToken()) != "" {
		rows, err = impl.pool.Query(ctx,
			"SELECT post_id, body FROM posts WHERE post_id < $1 ORDER BY post_id DESC LIMIT $2",
			string(*pageOption.PageToken()), pageSize,
		)
	} else {
		rows, err = impl.pool.Query(ctx,
			"SELECT post_id, body FROM posts ORDER BY post_id DESC LIMIT $1",
			pageSize,
		)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var postID, body string
		if err := rows.Scan(&postID, &body); err != nil {
			return nil, errors.WithStack(err)
		}
		posts = append(posts, entity.NewPost(entity.NewPostID(postID), entity.NewPostBody(body)))
	}
	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	if posts == nil {
		posts = []*entity.Post{}
	}
	return posts, nil
}
