package syncmap

import (
	"context"

	"github.com/becosuke/syncmap"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/adapter/repository"
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
)

func NewQuerier(config *pkgconfig.Config, logger *zap.Logger, store syncmap.Syncmap) repository.Querier {
	return &querierImpl{
		config: config,
		logger: logger,
		store:  store,
	}
}

type querierImpl struct {
	config *pkgconfig.Config
	logger *zap.Logger
	store  syncmap.Syncmap
}

func (impl *querierImpl) Get(_ context.Context, serial *domain.Serial) (*domain.Post, error) {
	value, err := impl.store.Get(serial.Int64())
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrInvalidArgument):
			return nil, repository.ErrInvalidArgument
		case errors.Is(err, syncmap.ErrNotFound):
			return nil, repository.ErrNotFound
		default:
			return nil, errors.WithStack(err)
		}
	}
	body, ok := value.(string)
	if !ok {
		return nil, repository.ErrInvalidData
	}

	return domain.NewPost(serial, domain.NewBody(body)), nil
}

func (impl *querierImpl) Range(_ context.Context, pageOption *domain.PageOption) ([]*domain.Post, error) {
	return []*domain.Post{}, nil
}
