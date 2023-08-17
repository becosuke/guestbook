package syncmap

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/drivers/syncmap"
	"github.com/becosuke/guestbook/api/internal/registry/config"
)

type Querier interface {
	Get(context.Context, *domain.Serial) (*domain.Post, error)
	Range(context.Context, *domain.PageOption) ([]*domain.Post, error)
}

func NewQuerier(config *config.Config, store syncmap.Syncmap, boundary Boundary) Querier {
	return &querierImpl{
		config:   config,
		store:    store,
		Boundary: boundary,
	}
}

type querierImpl struct {
	config *config.Config
	store  syncmap.Syncmap
	Boundary
}

func (impl *querierImpl) Get(ctx context.Context, serial *domain.Serial) (*domain.Post, error) {
	message, err := impl.store.Load(ctx, fmt.Sprintf("%d", serial.Int64()))
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrNotFound):
			return nil, ErrMessageNotFound
		case errors.Is(err, syncmap.ErrInvalidArgument):
			return nil, ErrInvalidArgument
		case errors.Is(err, syncmap.ErrInvalidData):
			return nil, ErrInvalidData
		default:
			return nil, errors.WithStack(err)
		}
	}
	return impl.ToEntity(message), nil
}

func (impl *querierImpl) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, error) {
	return []*domain.Post{}, nil
}
