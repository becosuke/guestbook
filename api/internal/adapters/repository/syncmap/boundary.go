package syncmap

import (
	"fmt"
	"strconv"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/drivers/syncmap"
)

type Boundary interface {
	ToEntity(message *syncmap.Message) *domain.Post
	ToMessage(entity *domain.Post) *syncmap.Message
}

func NewBoundary() Boundary {
	return &boundaryImpl{}
}

type boundaryImpl struct{}

func (impl *boundaryImpl) ToEntity(m *syncmap.Message) *domain.Post {
	if m == nil {
		return &domain.Post{}
	}
	serial, _ := strconv.ParseInt(m.Key(), 10, 64)
	return domain.NewPost(
		domain.NewSerial(serial),
		domain.NewBody(m.Value()),
	)
}

func (impl *boundaryImpl) ToMessage(post *domain.Post) *syncmap.Message {
	if post == nil {
		return &syncmap.Message{}
	}
	return syncmap.NewMessage(
		fmt.Sprintf("%d", post.Serial().Int64()),
		post.Body().String(),
	)
}
