package syncmap

import (
	"context"
	"sync"
)

type Syncmap interface {
	LoadOrStore(ctx context.Context, message *Message) (*Message, bool, error)
	Load(ctx context.Context, key string) (*Message, error)
	Store(ctx context.Context, message *Message) (*Message, error)
	Delete(ctx context.Context, key string) error
}

func NewSyncmap() Syncmap {
	return &syncmapImpl{
		syncmap: &sync.Map{},
	}
}

type syncmapImpl struct {
	syncmap *sync.Map
}

func (impl *syncmapImpl) LoadOrStore(_ context.Context, message *Message) (*Message, bool, error) {
	if message == nil {
		return nil, false, ErrSyncmapInvalidArgument
	}
	actual, loaded := impl.syncmap.LoadOrStore(message.Key(), message.Value())
	if loaded {
		asserted, ok := actual.(string)
		if !ok {
			return nil, loaded, ErrSyncmapInvalidData
		}
		return NewMessage(message.Key(), asserted), loaded, nil
	}
	return NewMessage(message.Key(), message.Value()), loaded, nil
}

func (impl *syncmapImpl) Load(_ context.Context, key string) (*Message, error) {
	if key == "" {
		return nil, ErrSyncmapInvalidArgument
	}
	actual, ok := impl.syncmap.Load(key)
	if !ok {
		return nil, ErrSyncmapNotFound
	}
	asserted, ok := actual.(string)
	if !ok {
		return nil, ErrSyncmapInvalidData
	}
	return NewMessage(key, asserted), nil
}

func (impl *syncmapImpl) Store(_ context.Context, message *Message) (*Message, error) {
	if message == nil {
		return nil, ErrSyncmapInvalidArgument
	}
	impl.syncmap.Store(message.Key(), message.Value())
	return NewMessage(message.Key(), message.Value()), nil
}

func (impl *syncmapImpl) Delete(_ context.Context, key string) error {
	if key == "" {
		return ErrSyncmapInvalidArgument
	}
	impl.syncmap.Delete(key)
	return nil
}
