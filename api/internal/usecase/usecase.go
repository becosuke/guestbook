package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

func NewUsecase(config *domain.Config, logger *zap.Logger, postQuerier interfaces.PostQuerier, postCommander interfaces.PostCommander, paginator interfaces.Paginator) *Usecase {
	return &Usecase{
		config:        config,
		logger:        logger,
		postQuerier:   postQuerier,
		postCommander: postCommander,
		paginator:     paginator,
	}
}

type Usecase struct {
	config        *domain.Config
	logger        *zap.Logger
	postQuerier   interfaces.PostQuerier
	postCommander interfaces.PostCommander
	paginator     interfaces.Paginator
}

func (impl *Usecase) Get(ctx context.Context, postID *domain.PostID) (*domain.Post, error) {
	result, err := impl.get(ctx, postID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (impl *Usecase) get(ctx context.Context, postID *domain.PostID) (*domain.Post, error) {
	result, err := impl.postQuerier.Get(ctx, postID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (impl *Usecase) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, *domain.PaginationID, error) {
	pageSize := int32(10)
	if pageOption.PageSize() != nil {
		pageSize = int32(*pageOption.PageSize())
	}

	var cursor *domain.PostCursor
	if pageOption.PageToken() != nil && string(*pageOption.PageToken()) != "" {
		pagination, err := impl.paginator.Get(ctx, domain.NewPaginationID(string(*pageOption.PageToken())))
		if err != nil {
			return nil, nil, err
		}
		c, err := domain.UnmarshalPostCursor(pagination.Cursor())
		if err != nil {
			return nil, nil, err
		}
		cursor = c
	}

	results, err := impl.postQuerier.Range(ctx, pageSize+1, cursor)
	if err != nil {
		return nil, nil, err
	}

	if int32(len(results)) > pageSize {
		lastPost := results[pageSize-1]
		results = results[:pageSize]

		nextCursor := domain.NewPostCursor(lastPost.PostID().String(), lastPost.CreateTime())
		cursorBytes, err := nextCursor.Marshal()
		if err != nil {
			return nil, nil, err
		}

		nextPaginationID := domain.NewPaginationID(uuid.New().String())
		pagination := domain.NewPagination(nextPaginationID, cursorBytes)
		if err := impl.paginator.Save(ctx, pagination); err != nil {
			return nil, nil, err
		}
		return results, nextPaginationID, nil
	}

	return results, nil, nil
}

func (impl *Usecase) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	postID := domain.NewPostID(uuid.New().String())
	post = domain.NewPost(postID, post.PostBody(), time.Time{}, nil)
	err := impl.postCommander.Create(ctx, post)
	if err != nil {
		return nil, err
	}
	return impl.get(ctx, postID)
}

func (impl *Usecase) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	err := impl.postCommander.Update(ctx, post)
	if err != nil {
		return nil, err
	}
	return impl.get(ctx, post.PostID())
}

func (impl *Usecase) Delete(ctx context.Context, postID *domain.PostID) error {
	err := impl.postCommander.Delete(ctx, postID)
	if err != nil {
		return err
	}
	return nil
}
