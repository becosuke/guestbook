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

func NewCommander(config *pkgconfig.Config, logger *zap.Logger, store syncmap.Syncmap, generator repository.Generator) repository.Commander {
	return &commanderImpl{
		config:    config,
		logger:    logger,
		store:     store,
		generator: generator,
	}
}

type commanderImpl struct {
	config    *pkgconfig.Config
	logger    *zap.Logger
	store     syncmap.Syncmap
	generator repository.Generator
}

func (impl *commanderImpl) Create(_ context.Context, post *domain.Post) (*domain.Serial, error) {
	serial := impl.generator.GenerateSerial()
	err := impl.store.Create(serial.Int64(), post.Body().String())
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrInvalidArgument):
			return nil, repository.ErrInvalidArgument
		case errors.Is(err, syncmap.ErrAlreadyExists):
			return nil, repository.ErrAlreadyExists
		default:
			return nil, errors.WithStack(err)
		}
	}
	return serial, nil
}

func (impl *commanderImpl) Update(_ context.Context, post *domain.Post) error {
	err := impl.store.Update(post.Serial().Int64(), post.Body().String())
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

func (impl *commanderImpl) Delete(_ context.Context, serial *domain.Serial) error {
	err := impl.store.Delete(serial.Int64())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
