package syncmap

import (
	"fmt"
	"strconv"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	driver "github.com/becosuke/guestbook/api/internal/driver/syncmap"
)

type Boundary interface {
	ToEntity(message *driver.Message) *domain.Post
	ToMessage(entity *domain.Post) *driver.Message
}

func NewBoundary() Boundary {
	return &boundaryImpl{}
}

type boundaryImpl struct{}

func (impl *boundaryImpl) ToEntity(m *driver.Message) *domain.Post {
	if m == nil {
		return &domain.Post{}
	}
	serial, _ := strconv.ParseInt(m.Key(), 10, 64)
	return domain.NewPost(
		domain.NewSerial(serial),
		domain.NewBody(m.Value()),
	)
}

func (impl *boundaryImpl) ToMessage(post *domain.Post) *driver.Message {
	if post == nil {
		return &driver.Message{}
	}
	return driver.NewMessage(
		fmt.Sprintf("%d", post.Serial().Int64()),
		post.Body().String(),
	)
}
