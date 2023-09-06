package controller

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	syncmap_repository "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

type guestbookServiceServerImpl struct {
	pb.UnimplementedGuestbookServiceServer
	config   *pkgconfig.Config
	logger   *zap.Logger
	usecase  usecase.Usecase
	boundary Boundary
}

func NewGuestbookServiceServer(config *pkgconfig.Config, logger *zap.Logger, usecase usecase.Usecase, boundary Boundary) pb.GuestbookServiceServer {
	return &guestbookServiceServerImpl{
		config:   config,
		logger:   logger,
		usecase:  usecase,
		boundary: boundary,
	}
}

func (impl *guestbookServiceServerImpl) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Get(ctx, impl.boundary.SerialResourceToDomain(req.GetSerial()))
	if err != nil {
		switch {
		case errors.Is(err, syncmap_repository.ErrMessageNotFound):
			return nil, status.New(codes.NotFound, err.Error()).Err()
		case errors.Is(err, syncmap_repository.ErrInvalidData), errors.Is(err, syncmap_repository.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.boundary.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Create(ctx, impl.boundary.PostResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, syncmap_repository.ErrInvalidData), errors.Is(err, syncmap_repository.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.boundary.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Update(ctx, impl.boundary.PostResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, syncmap_repository.ErrInvalidData), errors.Is(err, syncmap_repository.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.boundary.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*emptypb.Empty, error) {
	err := impl.usecase.Delete(ctx, impl.boundary.SerialResourceToDomain(req.GetSerial()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
