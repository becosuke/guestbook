package presentation

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/becosuke/guestbook/api/internal/adapter/presentation/converter"
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
	res, err := impl.usecase.Get(ctx, converter.PostIDResourceToDomain(req.GetPostId()))
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
	return converter.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	res, err := impl.usecase.Create(ctx, converter.PostResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidData), errors.Is(err, domain.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return converter.PostDomainToResource(res), nil
}

func (impl *guestbookServiceServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	if err := validateUpdateMask(req.GetUpdateMask()); err != nil {
		return nil, err
	}

	res, err := impl.usecase.Update(ctx, converter.PostResourceToDomain(req.GetPost()))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFailedPrecondition):
			return nil, status.New(codes.FailedPrecondition, err.Error()).Err()
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.New(codes.NotFound, err.Error()).Err()
		case errors.Is(err, domain.ErrInvalidData), errors.Is(err, domain.ErrInvalidArgument):
			return nil, status.New(codes.Internal, err.Error()).Err()
		default:
			return nil, status.New(codes.Unknown, err.Error()).Err()
		}
	}
	return converter.PostDomainToResource(res), nil
}

var validUpdateMaskPaths = map[string]bool{
	"body": true,
}

func validateUpdateMask(mask *fieldmaskpb.FieldMask) error {
	if mask == nil || len(mask.GetPaths()) == 0 {
		return nil
	}
	for _, path := range mask.GetPaths() {
		if !validUpdateMaskPaths[path] {
			return status.Errorf(codes.InvalidArgument, "invalid update_mask path: %q", path)
		}
	}
	return nil
}

func (impl *guestbookServiceServer) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	pageOption := converter.PageOptionResourceToDomain(req.GetPageSize(), req.GetPageToken())
	posts, nextPaginationID, err := impl.usecase.Range(ctx, pageOption)
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

	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, post := range posts {
		pbPosts = append(pbPosts, converter.PostDomainToResource(post))
	}

	resp := &pb.ListPostsResponse{
		Posts: pbPosts,
	}
	if !nextPaginationID.IsZero() {
		resp.NextPageToken = nextPaginationID.String()
	}
	return resp, nil
}

func (impl *guestbookServiceServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*emptypb.Empty, error) {
	err := impl.usecase.Delete(ctx, converter.PostIDResourceToDomain(req.GetPostId()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
