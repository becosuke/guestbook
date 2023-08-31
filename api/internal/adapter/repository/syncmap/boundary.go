package syncmap

import (
	"fmt"
	"strconv"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	syncmap_driver "github.com/becosuke/guestbook/api/internal/driver/syncmap"
)

type Boundary interface {
	ToEntity(message *syncmap_driver.Message) *domain.Post
	ToMessage(entity *domain.Post) *syncmap_driver.Message
}

func NewBoundary() Boundary {
	return &boundaryImpl{}
}

type boundaryImpl struct{}

func (impl *boundaryImpl) ToEntity(m *syncmap_driver.Message) *domain.Post {
	if m == nil {
		return &domain.Post{}
	}
	serial, _ := strconv.ParseInt(m.Key(), 10, 64)
	return domain.NewPost(
		domain.NewSerial(serial),
		domain.NewBody(m.Value()),
	)
}

func (impl *boundaryImpl) ToMessage(post *domain.Post) *syncmap_driver.Message {
	if post == nil {
		return &syncmap_driver.Message{}
	}
	return syncmap_driver.NewMessage(
		fmt.Sprintf("%d", post.Serial().Int64()),
		post.Body().String(),
	)
}
