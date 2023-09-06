package controller

import (
	"github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

type Boundary interface {
	SerialDomainToResource(domainSerial *post.Serial) int64
	SerialResourceToDomain(resourceSerial int64) *post.Serial
	BodyDomainToResource(domainBody *post.Body) string
	BodyResourceToDomain(resourceBody string) *post.Body
	PostDomainToResource(domainPost *post.Post) *pb.Post
	PostResourceToDomain(resourcePost *pb.Post) *post.Post
}

func NewBoundary() Boundary {
	return &boundaryImpl{}
}

type boundaryImpl struct{}

func (impl *boundaryImpl) SerialDomainToResource(domainSerial *post.Serial) int64 {
	return domainSerial.Int64()
}

func (impl *boundaryImpl) SerialResourceToDomain(resourceSerial int64) *post.Serial {
	return post.NewSerial(resourceSerial)
}

func (impl *boundaryImpl) BodyDomainToResource(domainBody *post.Body) string {
	return domainBody.String()
}

func (impl *boundaryImpl) BodyResourceToDomain(resourceBody string) *post.Body {
	return post.NewBody(resourceBody)
}

func (impl *boundaryImpl) PostDomainToResource(domainPost *post.Post) *pb.Post {
	return &pb.Post{
		Serial: impl.SerialDomainToResource(domainPost.Serial()),
		Body:   impl.BodyDomainToResource(domainPost.Body()),
	}
}

func (impl *boundaryImpl) PostResourceToDomain(resourcePost *pb.Post) *post.Post {
	return post.NewPost(
		impl.SerialResourceToDomain(resourcePost.GetSerial()),
		impl.BodyResourceToDomain(resourcePost.GetBody()),
	)
}
