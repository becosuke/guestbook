package usecase

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

func NewUsecase(config *domain.Config, logger *zap.Logger, repos interfaces.Repositories) *Usecase {
	return &Usecase{
		config: config,
		logger: logger,
		repos:  repos,
	}
}

type Usecase struct {
	config *domain.Config
	logger *zap.Logger
	repos  interfaces.Repositories
}

func (impl *Usecase) Get(ctx context.Context, postID domain.PostID) (*domain.Post, error) {
	result, err := impl.get(ctx, postID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (impl *Usecase) get(ctx context.Context, postID domain.PostID) (*domain.Post, error) {
	result, err := impl.repos.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// defaultPageSize は page_size が未指定（nil）または 0 のときに適用する件数。
// AIP-158 が認める「サーバ側で最終件数を決めてよい」枠を利用して、クライアントに
// デフォルトを暗黙適用する。
const defaultPageSize int32 = 10

func (impl *Usecase) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, domain.PaginationID, error) {
	pageSize := defaultPageSize
	if ps := pageOption.PageSize(); ps != nil && int32(*ps) > 0 {
		pageSize = int32(*ps)
	}

	var cursor *domain.PostCursor
	if pageOption.PageToken() != nil && string(*pageOption.PageToken()) != "" {
		pagination, err := impl.repos.GetPagination(ctx, domain.NewPaginationID(string(*pageOption.PageToken())))
		if err != nil {
			return nil, domain.PaginationID{}, err
		}
		c, err := domain.UnmarshalPostCursor(pagination.Cursor())
		if err != nil {
			return nil, domain.PaginationID{}, err
		}
		cursor = c
	}

	results, err := impl.repos.RangePosts(ctx, pageSize+1, cursor)
	if err != nil {
		return nil, domain.PaginationID{}, err
	}

	if int32(len(results)) > pageSize {
		lastPost := results[pageSize-1]
		results = results[:pageSize]

		nextCursor := domain.NewPostCursor(lastPost.PostID(), lastPost.CreateTime())
		cursorBytes, err := nextCursor.Marshal()
		if err != nil {
			return nil, domain.PaginationID{}, err
		}

		nextPaginationID := domain.NewPaginationID(uuid.NewString())
		pagination := domain.NewPagination(nextPaginationID, cursorBytes)
		if err := impl.repos.SavePagination(ctx, pagination); err != nil {
			return nil, domain.PaginationID{}, err
		}
		return results, nextPaginationID, nil
	}

	return results, domain.PaginationID{}, nil
}

func (impl *Usecase) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	newPost := domain.CreatePost(post.PostBody())
	if err := impl.repos.CreatePost(ctx, newPost); err != nil {
		return nil, err
	}
	return impl.get(ctx, newPost.PostID())
}

func (impl *Usecase) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	err := impl.repos.UpdatePost(ctx, post)
	if err != nil {
		return nil, err
	}
	return impl.get(ctx, post.PostID())
}

func (impl *Usecase) Delete(ctx context.Context, postID domain.PostID) error {
	err := impl.repos.DeletePost(ctx, postID)
	if err != nil {
		return err
	}
	return nil
}
