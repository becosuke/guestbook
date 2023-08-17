package syncmap

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/drivers/syncmap"
	"github.com/becosuke/guestbook/api/internal/registry/config"
)

type Commander interface {
	Create(context.Context, *domain.Post) (*domain.Serial, error)
	Update(context.Context, *domain.Post) error
	Delete(context.Context, *domain.Serial) error
}

func NewCommander(config *config.Config, store syncmap.Syncmap, boundary Boundary, generator Generator) Commander {
	return &commanderImpl{
		config:    config,
		store:     store,
		Boundary:  boundary,
		Generator: generator,
	}
}

type commanderImpl struct {
	config *config.Config
	store  syncmap.Syncmap
	Boundary
	Generator
}

func (impl *commanderImpl) Create(ctx context.Context, post *domain.Post) (*domain.Serial, error) {
	serial := impl.GenerateSerial()
	_, loaded, err := impl.store.LoadOrStore(ctx, impl.ToMessage(domain.NewPost(serial, post.Body())))
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrInvalidArgument):
			return nil, ErrInvalidArgument
		case errors.Is(err, syncmap.ErrInvalidData):
			return nil, ErrInvalidData
		default:
			return nil, errors.WithStack(err)
		}
	}
	if loaded {
		return nil, ErrMessageAlreadyExists
	}
	return serial, nil
}

func (impl *commanderImpl) Update(ctx context.Context, post *domain.Post) error {
	_, err := impl.store.Load(ctx, fmt.Sprintf("%d", post.Serial().Int64()))
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrNotFound):
			return ErrMessageNotFound
		case errors.Is(err, syncmap.ErrInvalidArgument):
			return ErrInvalidArgument
		case errors.Is(err, syncmap.ErrInvalidData):
			return ErrInvalidData
		}
	}
	_, err = impl.store.Store(ctx, impl.ToMessage(post))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (impl *commanderImpl) Delete(ctx context.Context, serial *domain.Serial) error {
	err := impl.store.Delete(ctx, fmt.Sprintf("%d", serial.Int64()))
	return errors.WithStack(err)
}
