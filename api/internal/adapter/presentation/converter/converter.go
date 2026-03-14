package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func PostIDDomainToResource(domainPostID domain.PostID) string {
	return domainPostID.String()
}

func PostIDResourceToDomain(resourcePostId string) domain.PostID {
	return domain.NewPostID(resourcePostId)
}

func PostBodyDomainToResource(domainPostBody domain.PostBody) string {
	return domainPostBody.String()
}

func PostBodyResourceToDomain(resourceBody string) domain.PostBody {
	return domain.NewPostBody(resourceBody)
}

func PostDomainToResource(domainPost *domain.Post) *pb.Post {
	if !domainPost.Valid() {
		return &pb.Post{
			PostId: PostIDDomainToResource(domainPost.PostID()),
			Valid:  false,
		}
	}
	res := &pb.Post{
		PostId:     PostIDDomainToResource(domainPost.PostID()),
		Body:       PostBodyDomainToResource(domainPost.PostBody()),
		Valid:      true,
		CreateTime: timestamppb.New(domainPost.CreateTime()),
	}
	if !domainPost.UpdateTime().IsZero() {
		res.UpdateTime = timestamppb.New(domainPost.UpdateTime())
	}
	if domainPost.PreviousBody().String() != "" {
		res.PreviousBody = PostBodyDomainToResource(domainPost.PreviousBody())
	}
	return res
}

func PostResourceToDomain(resourcePost *pb.Post) *domain.Post {
	return domain.NewPost(
		PostIDResourceToDomain(resourcePost.GetPostId()),
		PostBodyResourceToDomain(resourcePost.GetBody()),
		domain.NewPostBody(""),
		time.Time{},
		time.Time{},
		time.Time{},
	)
}

func PageOptionResourceToDomain(pageSize int32, pageToken string) *domain.PageOption {
	ps := domain.PageSize(pageSize)
	var pt *domain.PageToken
	if pageToken != "" {
		t := domain.PageToken(pageToken)
		pt = &t
	}
	return domain.NewPageOption(&ps, pt)
}
