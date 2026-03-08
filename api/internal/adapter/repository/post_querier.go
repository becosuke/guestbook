package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

func NewPostQuerier(config *domain.Config, logger *zap.Logger, pool *pgxpool.Pool) interfaces.PostQuerier {
	return &postQuerierImpl{
		config: config,
		logger: logger,
		pool:   pool,
	}
}

type postQuerierImpl struct {
	config *domain.Config
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func (impl *postQuerierImpl) GetPost(ctx context.Context, postID domain.PostID) (*domain.Post, error) {
	var body string
	var previousBody string
	var createTime time.Time
	var updateTime time.Time
	var deleteTime time.Time
	err := impl.pool.QueryRow(ctx, `SELECT PostBody, PreviousBody, CreateTime, UpdateTime, DeleteTime FROM Posts WHERE PostId = $1`, postID.String()).Scan(&body, &previousBody, &createTime, &updateTime, &deleteTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return domain.NewPost(postID, domain.NewPostBody(body), domain.NewPostBody(previousBody), createTime, updateTime, deleteTime), nil
}

func (impl *postQuerierImpl) RangePosts(ctx context.Context, pageSize int32, cursor *domain.PostCursor) ([]*domain.Post, error) {
	var rows pgx.Rows
	var err error
	if cursor != nil {
		rows, err = impl.pool.Query(ctx,
			`SELECT PostId, PostBody, PreviousBody, CreateTime, UpdateTime, DeleteTime FROM Posts WHERE (CreateTime < $1) OR (CreateTime = $1 AND PostId > $2) ORDER BY CreateTime DESC, PostId ASC LIMIT $3`,
			cursor.LastCreateTime, cursor.LastPostID, pageSize,
		)
	} else {
		rows, err = impl.pool.Query(ctx,
			`SELECT PostId, PostBody, PreviousBody, CreateTime, UpdateTime, DeleteTime FROM Posts ORDER BY CreateTime DESC, PostId ASC LIMIT $1`,
			pageSize,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var postID, body, previousBody string
		var createTime time.Time
		var updateTime time.Time
		var deleteTime time.Time
		if err := rows.Scan(&postID, &body, &previousBody, &createTime, &updateTime, &deleteTime); err != nil {
			return nil, err
		}
		posts = append(posts, domain.NewPost(domain.NewPostID(postID), domain.NewPostBody(body), domain.NewPostBody(previousBody), createTime, updateTime, deleteTime))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if posts == nil {
		posts = []*domain.Post{}
	}
	return posts, nil
}
