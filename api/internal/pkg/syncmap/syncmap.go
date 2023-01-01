package syncmap

import (
	"sync"
)

type Syncmap interface {
	LoadOrStore(message *Message) (*Message, bool, error)
	Load(serial int64) (*Message, error)
	Store(message *Message) (*Message, error)
	Delete(serial int64) error
}

func NewSyncmap() Syncmap {
	return &syncmapImpl{
		syncmap: &sync.Map{},
	}
}

type syncmapImpl struct {
	syncmap *sync.Map
}

func (impl *syncmapImpl) LoadOrStore(message *Message) (*Message, bool, error) {
	if message == nil {
		return nil, false, ErrSyncmapInvalidArgument
	}
	actual, loaded := impl.syncmap.LoadOrStore(message.Serial(), message.Body())
	if loaded {
		asserted, ok := actual.(string)
		if !ok {
			return nil, loaded, ErrSyncmapInvalidData
		}
		return NewMessage(message.Serial(), asserted), loaded, nil
	}
	return NewMessage(message.Serial(), message.Body()), loaded, nil
}

func (impl *syncmapImpl) Load(serial int64) (*Message, error) {
	if serial == 0 {
		return nil, ErrSyncmapInvalidArgument
	}
	actual, ok := impl.syncmap.Load(serial)
	if !ok {
		return nil, ErrSyncmapNotFound
	}
	asserted, ok := actual.(string)
	if !ok {
		return nil, ErrSyncmapInvalidData
	}
	return NewMessage(serial, asserted), nil
}

func (impl *syncmapImpl) Store(message *Message) (*Message, error) {
	if message == nil {
		return nil, ErrSyncmapInvalidArgument
	}
	impl.syncmap.Store(message.Serial(), message.Body())
	return NewMessage(message.Serial(), message.Body()), nil
}

func (impl *syncmapImpl) Delete(serial int64) error {
	if serial == 0 {
		return ErrSyncmapInvalidArgument
	}
	impl.syncmap.Delete(serial)
	return nil
}
