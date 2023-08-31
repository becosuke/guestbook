package usecase

import (
	"context"

	"github.com/pkg/errors"

	syncmap_repository "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
)

type Usecase interface {
	Get(context.Context, *domain.Serial) (*domain.Post, error)
	Range(context.Context, *domain.PageOption) ([]*domain.Post, error)
	Create(context.Context, *domain.Post) (*domain.Post, error)
	Update(context.Context, *domain.Post) (*domain.Post, error)
	Delete(context.Context, *domain.Serial) error
}

func NewUsecase(config *pkgconfig.Config, querier syncmap_repository.Querier, commander syncmap_repository.Commander) Usecase {
	return &usecaseImpl{
		config:    config,
		querier:   querier,
		commander: commander,
	}
}

type usecaseImpl struct {
	config    *pkgconfig.Config
	querier   syncmap_repository.Querier
	commander syncmap_repository.Commander
}

func (impl *usecaseImpl) Get(ctx context.Context, serial *domain.Serial) (*domain.Post, error) {
	result, err := impl.get(ctx, serial)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Add here any side effects of calling. For example, counting up the number of views.
	return result, nil
}

func (impl *usecaseImpl) get(ctx context.Context, serial *domain.Serial) (*domain.Post, error) {
	return impl.querier.Get(ctx, serial)
}

func (impl *usecaseImpl) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, error) {
	result, err := impl.querier.Range(ctx, pageOption)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (impl *usecaseImpl) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	serial, err := impl.commander.Create(ctx, post)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return impl.get(ctx, serial)
}

func (impl *usecaseImpl) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	err := impl.commander.Update(ctx, post)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return impl.get(ctx, post.Serial())
}

func (impl *usecaseImpl) Delete(ctx context.Context, serial *domain.Serial) error {
	err := impl.commander.Delete(ctx, serial)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
