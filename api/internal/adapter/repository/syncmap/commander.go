package syncmap

import (
	"context"

	"github.com/becosuke/syncmap"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain/repository"
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
)

func NewCommander(config *pkgconfig.Config, logger *zap.Logger, store syncmap.Syncmap) repository.Commander {
	return &commanderImpl{
		config: config,
		logger: logger,
		store:  store,
	}
}

type commanderImpl struct {
	config *pkgconfig.Config
	logger *zap.Logger
	store  syncmap.Syncmap
}

func (impl *commanderImpl) Create(_ context.Context, post *domain.Post) error {
	err := impl.store.Create(post.PostID().String(), post.Body().String())
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrInvalidArgument):
			return repository.ErrInvalidArgument
		case errors.Is(err, syncmap.ErrAlreadyExists):
			return repository.ErrAlreadyExists
		default:
			return errors.WithStack(err)
		}
	}
	return nil
}

func (impl *commanderImpl) Update(_ context.Context, post *domain.Post) error {
	err := impl.store.Update(post.PostID().String(), post.Body().String())
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrInvalidArgument):
			return repository.ErrInvalidArgument
		case errors.Is(err, syncmap.ErrNotFound):
			return repository.ErrNotFound
		default:
			return errors.WithStack(err)
		}
	}
	return nil
}

func (impl *commanderImpl) Delete(_ context.Context, postID *domain.PostID) error {
	err := impl.store.Delete(postID.String())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
