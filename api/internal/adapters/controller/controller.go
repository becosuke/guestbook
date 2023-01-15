package controller

import (
	"context"
	"github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/pkg/syncmap"
	"github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/api/pb"
	"github.com/mennanov/fmutils"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type guestbookServiceServerImpl struct {
	pb.UnimplementedGuestbookServiceServer
	config   *config.Config
	usecase  post.Usecase
	boundary Boundary
}

func NewGuestbookServiceServer(config *config.Config, usecase post.Usecase, boundary Boundary) pb.GuestbookServiceServer {
	return &guestbookServiceServerImpl{
		config:   config,
		usecase:  usecase,
		boundary: boundary,
	}
}

func (impl *guestbookServiceServerImpl) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	res, err := impl.usecase.Get(ctx, impl.boundary.SerialResourceToDomain(req.GetSerial()))
	if err != nil {
		switch {
		case errors.Is(err, syncmap.ErrSyncmapInvalidArgument):
		case errors.Is(err, syncmap.ErrSyncmapNotFound):
		case errors.Is(err, syncmap.ErrSyncmapInvalidData):
		default:

		}
	}
	return impl.boundary.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	res, err := impl.usecase.Create(ctx, impl.boundary.PostResourceToDomain(req.GetPost()))
	if err != nil {
		return nil, err
	}
	return impl.boundary.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	dest := req.GetPost()
	req.GetFieldMask().Normalize()
	if req.GetFieldMask().IsValid(req.GetPost()) {
		fmutils.Filter(dest, req.GetFieldMask().GetPaths())
	}
	res, err := impl.usecase.Update(ctx, impl.boundary.PostResourceToDomain(dest))
	if err != nil {
		return nil, err
	}
	return impl.boundary.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	err := impl.usecase.Delete(ctx, impl.boundary.SerialResourceToDomain(req.GetSerial()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
