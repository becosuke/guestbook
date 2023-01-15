package gateway

import (
	"context"
	"fmt"
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/pkg/syncmap"
	"github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/pkg/errors"
	"strconv"
)

var (
	ErrRepositoryAlreadyExists   = errors.New("message already exists")
	ErrRepositoryNotFound        = errors.New("not exists")
	ErrRepositoryInvalidData     = errors.New("returns invalid data")
	ErrRepositoryInvalidArgument = errors.New("invalid argument")
)

func NewRepository(config *config.Config, syncmap syncmap.Syncmap, generator domain.Generator) domain.Repository {
	return &repositoryImpl{
		config:    config,
		store:     syncmap,
		generator: generator,
	}
}

type repositoryImpl struct {
	config    *config.Config
	store     syncmap.Syncmap
	generator domain.Generator
}

func (impl *repositoryImpl) Get(ctx context.Context, serial *domain.Serial) (*domain.Post, error) {
	message, err := impl.store.Load(fmt.Sprintf("%d", serial.Int64()))
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrSyncmapNotFound):
			return nil, ErrRepositoryNotFound
		case errors.Is(err, syncmap.ErrSyncmapInvalidArgument):
			return nil, ErrRepositoryInvalidArgument
		case errors.Is(err, syncmap.ErrSyncmapInvalidData):
			return nil, ErrRepositoryInvalidData
		default:
			return nil, errors.WithStack(err)
		}
	}
	return impl.ToEntity(message), nil
}

func (impl *repositoryImpl) Range(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, error) {
	return []*domain.Post{}, nil
}

func (impl *repositoryImpl) Create(ctx context.Context, post *domain.Post) (*domain.Serial, error) {
	serial := impl.generator.GenerateSerial()
	_, loaded, err := impl.store.LoadOrStore(impl.ToMessage(domain.NewPost(serial, post.Body())))
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrSyncmapInvalidArgument):
			return nil, ErrRepositoryInvalidArgument
		case errors.Is(err, syncmap.ErrSyncmapInvalidData):
			return nil, ErrRepositoryInvalidData
		default:
			return nil, errors.WithStack(err)
		}
	}
	if loaded {
		return nil, ErrRepositoryAlreadyExists
	}
	return serial, nil
}

func (impl *repositoryImpl) Update(ctx context.Context, post *domain.Post) error {
	_, err := impl.store.Load(fmt.Sprintf("%d", post.Serial().Int64()))
	if err == nil {
		switch {
		case errors.Is(err, syncmap.ErrSyncmapNotFound):
			return ErrRepositoryNotFound
		case errors.Is(err, syncmap.ErrSyncmapInvalidArgument):
			return ErrRepositoryInvalidArgument
		case errors.Is(err, syncmap.ErrSyncmapInvalidData):
			return ErrRepositoryInvalidData
		}
	}
	_, err = impl.store.Store(impl.ToMessage(post))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (impl *repositoryImpl) Delete(ctx context.Context, serial *domain.Serial) error {
	err := impl.store.Delete(fmt.Sprintf("%d", serial.Int64()))
	return errors.WithStack(err)
}

func (impl *repositoryImpl) ToEntity(m *syncmap.Message) *domain.Post {
	if m == nil {
		return &domain.Post{}
	}
	serial, _ := strconv.ParseInt(m.Key(), 10, 64)
	return domain.NewPost(
		domain.NewSerial(serial),
		domain.NewBody(m.Value()),
	)
}

func (impl *repositoryImpl) ToMessage(post *domain.Post) *syncmap.Message {
	if post == nil {
		return &syncmap.Message{}
	}
	return syncmap.NewMessage(
		fmt.Sprintf("%d", post.Serial().Int64()),
		post.Body().String(),
	)
}
