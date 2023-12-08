package controller

import (
	"github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func (impl *guestbookServiceServerImpl) serialDomainToResource(domainSerial *post.Serial) int64 {
	return domainSerial.Int64()
}

func (impl *guestbookServiceServerImpl) serialResourceToDomain(resourceSerial int64) *post.Serial {
	return post.NewSerial(resourceSerial)
}

func (impl *guestbookServiceServerImpl) bodyDomainToResource(domainBody *post.Body) string {
	return domainBody.String()
}

func (impl *guestbookServiceServerImpl) bodyResourceToDomain(resourceBody string) *post.Body {
	return post.NewBody(resourceBody)
}

func (impl *guestbookServiceServerImpl) postDomainToResource(domainPost *post.Post) *pb.Post {
	return &pb.Post{
		Serial: impl.serialDomainToResource(domainPost.Serial()),
		Body:   impl.bodyDomainToResource(domainPost.Body()),
	}
}

func (impl *guestbookServiceServerImpl) postResourceToDomain(resourcePost *pb.Post) *post.Post {
	return post.NewPost(
		impl.serialResourceToDomain(resourcePost.GetSerial()),
		impl.bodyResourceToDomain(resourcePost.GetBody()),
	)
}
