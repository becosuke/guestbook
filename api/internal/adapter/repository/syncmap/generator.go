package syncmap

import (
	"sync/atomic"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
)

type Generator interface {
	GenerateSerial() *domain.Serial
}

func NewGenerator() Generator {
	return &generatorImpl{}
}

type generatorImpl struct {
	serial int64
}

func (impl *generatorImpl) GenerateSerial() *domain.Serial {
	return domain.NewSerial(atomic.AddInt64(&impl.serial, 1))
}
