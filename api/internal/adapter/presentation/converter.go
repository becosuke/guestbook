package presentation

import (
	"github.com/becosuke/guestbook/api/internal/domain/entity"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func (impl *guestbookServiceServer) postIDDomainToResource(domainPostID *entity.PostID) string {
	return domainPostID.String()
}

func (impl *guestbookServiceServer) postIDResourceToDomain(resourcePostId string) *entity.PostID {
	return entity.NewPostID(resourcePostId)
}

func (impl *guestbookServiceServer) postBodyDomainToResource(domainPostBody *entity.PostBody) string {
	return domainPostBody.String()
}

func (impl *guestbookServiceServer) postBodyResourceToDomain(resourceBody string) *entity.PostBody {
	return entity.NewPostBody(resourceBody)
}

func (impl *guestbookServiceServer) postDomainToResource(domainPost *entity.Post) *pb.Post {
	return &pb.Post{
		PostId: impl.postIDDomainToResource(domainPost.PostID()),
		Body:   impl.postBodyDomainToResource(domainPost.PostBody()),
	}
}

func (impl *guestbookServiceServer) postResourceToDomain(resourcePost *pb.Post) *entity.Post {
	return entity.NewPost(
		impl.postIDResourceToDomain(resourcePost.GetPostId()),
		impl.postBodyResourceToDomain(resourcePost.GetBody()),
	)
}
