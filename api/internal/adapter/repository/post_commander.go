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

func NewPostCommander(config *domain.Config, logger *zap.Logger, pool *pgxpool.Pool) interfaces.PostCommander {
	return &postCommanderImpl{
		config: config,
		logger: logger,
		pool:   pool,
	}
}

type postCommanderImpl struct {
	config *domain.Config
	logger *zap.Logger
	pool   *pgxpool.Pool
}

func (impl *postCommanderImpl) CreatePost(ctx context.Context, post *domain.Post) error {
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

// UpdatePost は投稿の本文を 1 度だけ書き換える。
//
// WHERE 句の `CreateTime = UpdateTime` は「まだ一度も書き直されていない」ことを示す
// 前提条件として使う。INSERT 時は CreateTime / UpdateTime がどちらも NOW() で同値、
// 初回の UPDATE で UpdateTime のみが NOW() に書き換わるため、以後この等値は崩れる。
// 結果として 2 度目以降の UpdatePost は対象行に当たらず、1 投稿につき書き直しは
// 1 回までという仕様を SQL レベルで担保する。
//
// 0 行更新になった場合の振り分け:
//   - レコードが存在する  → 削除済み or 既に書き直し済み → ErrFailedPrecondition
//   - レコードが存在しない → 該当 post_id 自体が無い         → ErrNotFound
//
// 「削除済み」と「既に書き直し済み」を ErrFailedPrecondition に集約しているため、
// 両者を区別して通知したい場合は presentation 層で google.rpc.PreconditionFailure
// （google.rpc.ErrorDetails の一種）を付与する必要がある。現状はそこまで使い分けず、
// 識別性を犠牲にして単一の sentinel error にまとめている。
func (impl *postCommanderImpl) UpdatePost(ctx context.Context, post *domain.Post) error {
	ct, err := impl.pool.Exec(ctx,
		`UPDATE Posts SET PreviousBody = PostBody, PostBody = $1, UpdateTime = NOW() WHERE PostId = $2 AND DeleteTime = '0001-01-01 00:00:00+00' AND CreateTime = UpdateTime`,
		post.PostBody().String(), post.PostID().String(),
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		var exists bool
		err := impl.pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM Posts WHERE PostId = $1)`,
			post.PostID().String(),
		).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrFailedPrecondition
		}
		return domain.ErrNotFound
	}
	return nil
}

func (impl *postCommanderImpl) DeletePost(ctx context.Context, postID domain.PostID) error {
	ct, err := impl.pool.Exec(ctx,
		`UPDATE Posts SET DeleteTime = NOW(), UpdateTime = NOW() WHERE PostId = $1 AND DeleteTime = '0001-01-01 00:00:00+00'`,
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
