package presentation

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func (impl *guestbookServiceServer) postIDDomainToResource(domainPostID domain.PostID) string {
	return domainPostID.String()
}

func (impl *guestbookServiceServer) postIDResourceToDomain(resourcePostId string) domain.PostID {
	return domain.NewPostID(resourcePostId)
}

func (impl *guestbookServiceServer) postBodyDomainToResource(domainPostBody domain.PostBody) string {
	return domainPostBody.String()
}

func (impl *guestbookServiceServer) postBodyResourceToDomain(resourceBody string) domain.PostBody {
	return domain.NewPostBody(resourceBody)
}

func (impl *guestbookServiceServer) postDomainToResource(domainPost *domain.Post) *pb.Post {
	if !domainPost.Valid() {
		return &pb.Post{
			PostId: impl.postIDDomainToResource(domainPost.PostID()),
			Valid:  false,
		}
	}
	res := &pb.Post{
		PostId:     impl.postIDDomainToResource(domainPost.PostID()),
		Body:       impl.postBodyDomainToResource(domainPost.PostBody()),
		Valid:      true,
		CreateTime: timestamppb.New(domainPost.CreateTime()),
	}
	if !domainPost.UpdateTime().IsZero() {
		res.UpdateTime = timestamppb.New(domainPost.UpdateTime())
	}
	if domainPost.PreviousBody().String() != "" {
		res.PreviousBody = impl.postBodyDomainToResource(domainPost.PreviousBody())
	}
	return res
}

func (impl *guestbookServiceServer) postResourceToDomain(resourcePost *pb.Post) *domain.Post {
	return domain.NewPost(
		impl.postIDResourceToDomain(resourcePost.GetPostId()),
		impl.postBodyResourceToDomain(resourcePost.GetBody()),
		domain.NewPostBody(""),
		time.Time{},
		time.Time{},
		time.Time{},
	)
}

func (impl *guestbookServiceServer) pageOptionResourceToDomain(pageSize int32, pageToken string) *domain.PageOption {
	ps := domain.PageSize(pageSize)
	var pt *domain.PageToken
	if pageToken != "" {
		t := domain.PageToken(pageToken)
		pt = &t
	}
	return domain.NewPageOption(&ps, pt)
}
