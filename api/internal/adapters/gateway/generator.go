package gateway

import (
	"github.com/becosuke/guestbook/api/internal/domain/post"
)

func NewGenerator() post.Generator {
	return &generatorImpl{}
}

type generatorImpl struct{}

func (impl *generatorImpl) GenerateSerial() *post.Serial {
	s := post.Serial(1)
	return &s
}
