package presentation

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

type guestbookServiceServer struct {
	pb.UnimplementedGuestbookServiceServer
	config  *domain.Config
	logger  *zap.Logger
	usecase Usecase
}

func NewGuestbookServiceServer(config *domain.Config, logger *zap.Logger, usecase Usecase) pb.GuestbookServiceServer {
	return &guestbookServiceServer{
		config:  config,
		logger:  logger,
		usecase: usecase,
	}
}

func (impl *guestbookServiceServer) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Get(ctx, impl.postIDResourceToDomain(req.GetPostId()))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.New(codes.NotFound, err.Error()).Err()
		case errors.Is(err, domain.ErrInvalidData), errors.Is(err, domain.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.postDomainToResource(res), nil
}

func (impl *guestbookServiceServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Create(ctx, impl.postResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidData), errors.Is(err, domain.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.postDomainToResource(res), nil
}

func (impl *guestbookServiceServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Update(ctx, impl.postResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidData), errors.Is(err, domain.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return impl.postDomainToResource(res), nil
}

func (impl *guestbookServiceServer) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	return &pb.ListPostsResponse{}, nil
}

func (impl *guestbookServiceServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*emptypb.Empty, error) {
	err := impl.usecase.Delete(ctx, impl.postIDResourceToDomain(req.GetPostId()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
