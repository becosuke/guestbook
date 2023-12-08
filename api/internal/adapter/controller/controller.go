package controller

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/becosuke/guestbook/api/internal/adapter/repository"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

type guestbookServiceServerImpl struct {
	pb.UnimplementedGuestbookServiceServer
	config  *pkgconfig.Config
	logger  *zap.Logger
	usecase usecase.Usecase
}

func NewGuestbookServiceServer(config *pkgconfig.Config, logger *zap.Logger, usecase usecase.Usecase) pb.GuestbookServiceServer {
	return &guestbookServiceServerImpl{
		config:  config,
		logger:  logger,
		usecase: usecase,
	}
}

func (impl *guestbookServiceServerImpl) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Get(ctx, impl.serialResourceToDomain(req.GetSerial()))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, status.New(codes.NotFound, err.Error()).Err()
		case errors.Is(err, repository.ErrInvalidData), errors.Is(err, repository.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.postDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Create(ctx, impl.postResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidData), errors.Is(err, repository.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.postDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Update(ctx, impl.postResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidData), errors.Is(err, repository.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.postDomainToResource(res), nil
}

func (impl *guestbookServiceServerImpl) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*emptypb.Empty, error) {
	err := impl.usecase.Delete(ctx, impl.serialResourceToDomain(req.GetSerial()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
