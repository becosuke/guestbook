package controller

import (
	"context"

	"github.com/mennanov/fmutils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	repository "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	"github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/pbgo"
)

type guestbookServiceServerImpl struct {
	pbgo.UnimplementedGuestbookServiceServer
	config   *config.Config
	usecase  usecase.Usecase
	boundary Boundary
}

func NewGuestbookServiceServer(config *config.Config, usecase usecase.Usecase, boundary Boundary) pbgo.GuestbookServiceServer {
	return &guestbookServiceServerImpl{
		config:   config,
		usecase:  usecase,
		boundary: boundary,
	}
}

func (impl *guestbookServiceServerImpl) GetPost(ctx context.Context, req *pbgo.GetPostRequest) (*pbgo.Post, error) {
	res, err := impl.usecase.Get(ctx, impl.boundary.SerialResourceToDomain(req.GetSerial()))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrMessageNotFound):
			return nil, status.New(codes.NotFound, err.Error()).Err()
		case errors.Is(err, repository.ErrInvalidData), errors.Is(err, repository.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.boundary.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) CreatePost(ctx context.Context, req *pbgo.CreatePostRequest) (*pbgo.Post, error) {
	res, err := impl.usecase.Create(ctx, impl.boundary.PostResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidData), errors.Is(err, repository.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.boundary.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) UpdatePost(ctx context.Context, req *pbgo.UpdatePostRequest) (*pbgo.Post, error) {
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

func (impl *guestbookServiceServerImpl) DeletePost(ctx context.Context, req *pbgo.DeletePostRequest) (*emptypb.Empty, error) {
	err := impl.usecase.Delete(ctx, impl.boundary.SerialResourceToDomain(req.GetSerial()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
