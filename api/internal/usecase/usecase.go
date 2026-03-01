package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

func NewUsecase(config *domain.Config, logger *zap.Logger, querier interfaces.Querier, commander interfaces.Commander) *Usecase {
	return &Usecase{
		config:    config,
		logger:    logger,
		querier:   querier,
		commander: commander,
	}
}

type Usecase struct {
	config    *domain.Config
	logger    *zap.Logger
	querier   interfaces.Querier
	commander interfaces.Commander
}

func (impl *Usecase) Get(ctx context.Context, postID *domain.PostID) (*domain.Post, error) {
	result, err := impl.get(ctx, postID)
	if err != nil {
		return nil, err // Already stacked
	}
	// Add here any side effects of calling. For example, counting up the number of views.
	return result, nil
}

func (impl *Usecase) get(ctx context.Context, postID *domain.PostID) (*domain.Post, error) {
	result, err := impl.querier.Get(ctx, postID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (impl *Usecase) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, error) {
	result, err := impl.querier.Range(ctx, pageOption)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (impl *Usecase) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	postID := domain.NewPostID(uuid.New().String())
	post = domain.NewPost(postID, post.PostBody())
	err := impl.commander.Create(ctx, post)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return impl.get(ctx, postID)
}

func (impl *Usecase) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	err := impl.commander.Update(ctx, post)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return impl.get(ctx, post.PostID())
}

func (impl *Usecase) Delete(ctx context.Context, postID *domain.PostID) error {
	err := impl.commander.Delete(ctx, postID)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
