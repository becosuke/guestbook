package presentation

import (
	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func (impl *guestbookServiceServer) postIDDomainToResource(domainPostID *domain.PostID) string {
	return domainPostID.String()
}

func (impl *guestbookServiceServer) postIDResourceToDomain(resourcePostId string) *domain.PostID {
	return domain.NewPostID(resourcePostId)
}

func (impl *guestbookServiceServer) postBodyDomainToResource(domainPostBody *domain.PostBody) string {
	return domainPostBody.String()
}

func (impl *guestbookServiceServer) postBodyResourceToDomain(resourceBody string) *domain.PostBody {
	return domain.NewPostBody(resourceBody)
}

func (impl *guestbookServiceServer) postDomainToResource(domainPost *domain.Post) *pb.Post {
	if !domainPost.Valid() {
		return &pb.Post{
			PostId: impl.postIDDomainToResource(domainPost.PostID()),
			Valid:  false,
		}
	}
	return &pb.Post{
		PostId: impl.postIDDomainToResource(domainPost.PostID()),
		Body:   impl.postBodyDomainToResource(domainPost.PostBody()),
		Valid:  true,
	}
}

func (impl *guestbookServiceServer) postResourceToDomain(resourcePost *pb.Post) *domain.Post {
	return domain.NewPost(
		impl.postIDResourceToDomain(resourcePost.GetPostId()),
		impl.postBodyResourceToDomain(resourcePost.GetBody()),
		nil,
	)
}
