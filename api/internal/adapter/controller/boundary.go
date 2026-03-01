package controller

import (
	"strconv"

	"github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func (impl *guestbookServiceServerImpl) serialDomainToResource(domainSerial *post.Serial) string {
	return strconv.FormatInt(domainSerial.Int64(), 10)
}

func (impl *guestbookServiceServerImpl) serialResourceToDomain(resourcePostId string) *post.Serial {
	v, _ := strconv.ParseInt(resourcePostId, 10, 64)
	return post.NewSerial(v)
}

func (impl *guestbookServiceServerImpl) bodyDomainToResource(domainBody *post.Body) string {
	return domainBody.String()
}

func (impl *guestbookServiceServerImpl) bodyResourceToDomain(resourceBody string) *post.Body {
	return post.NewBody(resourceBody)
}

func (impl *guestbookServiceServerImpl) postDomainToResource(domainPost *post.Post) *pb.Post {
	return &pb.Post{
		PostId: impl.serialDomainToResource(domainPost.Serial()),
		Body:   impl.bodyDomainToResource(domainPost.Body()),
	}
}

func (impl *guestbookServiceServerImpl) postResourceToDomain(resourcePost *pb.Post) *post.Post {
	return post.NewPost(
		impl.serialResourceToDomain(resourcePost.GetPostId()),
		impl.bodyResourceToDomain(resourcePost.GetBody()),
	)
}
