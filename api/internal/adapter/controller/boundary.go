package controller

import (
	"github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func (impl *guestbookServiceServerImpl) postIDDomainToResource(domainPostID *post.PostID) string {
	return domainPostID.String()
}

func (impl *guestbookServiceServerImpl) postIDResourceToDomain(resourcePostId string) *post.PostID {
	return post.NewPostID(resourcePostId)
}

func (impl *guestbookServiceServerImpl) bodyDomainToResource(domainBody *post.Body) string {
	return domainBody.String()
}

func (impl *guestbookServiceServerImpl) bodyResourceToDomain(resourceBody string) *post.Body {
	return post.NewBody(resourceBody)
}

func (impl *guestbookServiceServerImpl) postDomainToResource(domainPost *post.Post) *pb.Post {
	return &pb.Post{
		PostId: impl.postIDDomainToResource(domainPost.PostID()),
		Body:   impl.bodyDomainToResource(domainPost.Body()),
	}
}

func (impl *guestbookServiceServerImpl) postResourceToDomain(resourcePost *pb.Post) *post.Post {
	return post.NewPost(
		impl.postIDResourceToDomain(resourcePost.GetPostId()),
		impl.bodyResourceToDomain(resourcePost.GetBody()),
	)
}
