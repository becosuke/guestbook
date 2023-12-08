package generator

import (
	"sync/atomic"

	"github.com/becosuke/guestbook/api/internal/adapter/repository"
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
)

func NewGenerator() repository.Generator {
	return &generatorImpl{}
}

type generatorImpl struct {
	serial int64
}

func (impl *generatorImpl) GenerateSerial() *domain.Serial {
	return domain.NewSerial(atomic.AddInt64(&impl.serial, 1))
}
