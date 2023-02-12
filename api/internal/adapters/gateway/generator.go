package gateway

import (
	"github.com/becosuke/guestbook/api/internal/domain/post"
	"sync/atomic"
)

func NewGenerator() post.Generator {
	return &generatorImpl{}
}

type generatorImpl struct {
	serial int64
}

func (impl *generatorImpl) GenerateSerial() *post.Serial {
	return post.NewSerial(atomic.AddInt64(&impl.serial, 1))
}
