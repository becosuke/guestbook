package interactor

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/adapter/repository"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
)

func NewUsecase(config *pkgconfig.Config, logger *zap.Logger, querier repository.Querier, commander repository.Commander) usecase.Usecase {
	return &usecaseImpl{
		config:    config,
		logger:    logger,
		querier:   querier,
		commander: commander,
	}
}

type usecaseImpl struct {
	config    *pkgconfig.Config
	logger    *zap.Logger
	querier   repository.Querier
	commander repository.Commander
}

func (impl *usecaseImpl) Get(ctx context.Context, postID *domain.PostID) (*domain.Post, error) {
	result, err := impl.get(ctx, postID)
	if err != nil {
		return nil, err // Already stacked
	}
	// Add here any side effects of calling. For example, counting up the number of views.
	return result, nil
}

func (impl *usecaseImpl) get(ctx context.Context, postID *domain.PostID) (*domain.Post, error) {
	result, err := impl.querier.Get(ctx, postID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (impl *usecaseImpl) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, error) {
	result, err := impl.querier.Range(ctx, pageOption)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (impl *usecaseImpl) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	postID := domain.NewPostID(uuid.New().String())
	post = domain.NewPost(postID, post.Body())
	err := impl.commander.Create(ctx, post)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return impl.get(ctx, postID)
}

func (impl *usecaseImpl) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	err := impl.commander.Update(ctx, post)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return impl.get(ctx, post.PostID())
}

func (impl *usecaseImpl) Delete(ctx context.Context, postID *domain.PostID) error {
	err := impl.commander.Delete(ctx, postID)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
