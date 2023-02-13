package usecase

import (
	"context"
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/pkg/errors"
)

func NewUsecase(config *config.Config, repository domain.Repository) domain.Usecase {
	return &usecaseImpl{
		config:     config,
		repository: repository,
	}
}

type usecaseImpl struct {
	config     *config.Config
	repository domain.Repository
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
	return impl.repository.Get(ctx, serial)
}

func (impl *usecaseImpl) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, error) {
	result, err := impl.repository.Range(ctx, pageOption)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (impl *usecaseImpl) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	serial, err := impl.repository.Create(ctx, post)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return impl.get(ctx, serial)
}

func (impl *usecaseImpl) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	err := impl.repository.Update(ctx, post)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return impl.get(ctx, post.Serial())
}

func (impl *usecaseImpl) Delete(ctx context.Context, serial *domain.Serial) error {
	err := impl.repository.Delete(ctx, serial)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
