package presentation

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/becosuke/guestbook/api/internal/domain"
	pb "github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func Test_guestbookServiceServerImpl_GetPost(t *testing.T) {
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *domain.Config
		usecase                             Usecase
	}
	type args struct {
		ctx context.Context
		req *pb.GetPostRequest
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *pb.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			config := &domain.Config{}
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := domain.NewPostBody("example")
			post := domain.NewPost(postID, body, domain.NewPostBody(""), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}, time.Time{})
			mockUsecase := &UsecaseMock{
				GetFunc: func(ctx context.Context, id domain.PostID) (*domain.Post, error) {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
					return post, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					config:  config,
					usecase: mockUsecase,
				},
				args: args{
					ctx: ctx,
					req: &pb.GetPostRequest{
						PostId: "550e8400-e29b-41d4-a716-446655440000",
					},
				},
				want: &pb.Post{
					PostId:     "550e8400-e29b-41d4-a716-446655440000",
					Body:       "example",
					Valid:      true,
					CreateTime: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			config := &domain.Config{}
			mockUsecase := &UsecaseMock{
				GetFunc: func(ctx context.Context, id domain.PostID) (*domain.Post, error) {
					return nil, domain.ErrNotFound
				},
			}
			return testCase{
				name: "not found",
				fields: fields{
					config:  config,
					usecase: mockUsecase,
				},
				args: args{
					ctx: ctx,
					req: &pb.GetPostRequest{
						PostId: "550e8400-e29b-41d4-a716-446655440000",
					},
				},
				want:    nil,
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServer{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
			}
			got, err := impl.GetPost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServer.GetPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServer.GetPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guestbookServiceServerImpl_CreatePost(t *testing.T) {
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *domain.Config
		usecase                             Usecase
	}
	type args struct {
		ctx context.Context
		req *pb.CreatePostRequest
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *pb.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			config := &domain.Config{}
			res := domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example"), domain.NewPostBody(""), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}, time.Time{})
			mockUsecase := &UsecaseMock{
				CreateFunc: func(ctx context.Context, post *domain.Post) (*domain.Post, error) {
					return res, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					config:  config,
					usecase: mockUsecase,
				},
				args: args{
					ctx: ctx,
					req: &pb.CreatePostRequest{
						Post: &pb.Post{
							Body: "example",
						},
					},
				},
				want: &pb.Post{
					PostId:     "550e8400-e29b-41d4-a716-446655440000",
					Body:       "example",
					Valid:      true,
					CreateTime: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServer{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
			}
			got, err := impl.CreatePost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServer.CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServer.CreatePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guestbookServiceServerImpl_UpdatePost(t *testing.T) {
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *domain.Config
		usecase                             Usecase
	}
	type args struct {
		ctx context.Context
		req *pb.UpdatePostRequest
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *pb.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			config := &domain.Config{}
			post := domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example-value"), domain.NewPostBody(""), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}, time.Time{})
			mockUsecase := &UsecaseMock{
				UpdateFunc: func(ctx context.Context, p *domain.Post) (*domain.Post, error) {
					return post, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					config:  config,
					usecase: mockUsecase,
				},
				args: args{
					ctx: ctx,
					req: &pb.UpdatePostRequest{
						Post: &pb.Post{
							PostId: "550e8400-e29b-41d4-a716-446655440000",
							Body:   "example-value",
						},
					},
				},
				want: &pb.Post{
					PostId:     "550e8400-e29b-41d4-a716-446655440000",
					Body:       "example-value",
					Valid:      true,
					CreateTime: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServer{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
			}
			got, err := impl.UpdatePost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServer.UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServer.UpdatePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guestbookServiceServerImpl_DeletePost(t *testing.T) {
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *domain.Config
		usecase                             Usecase
	}
	type args struct {
		ctx context.Context
		req *pb.DeletePostRequest
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			config := &domain.Config{}
			mockUsecase := &UsecaseMock{
				DeleteFunc: func(ctx context.Context, id domain.PostID) error {
					return nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					config:  config,
					usecase: mockUsecase,
				},
				args: args{
					ctx: ctx,
					req: &pb.DeletePostRequest{
						PostId: "550e8400-e29b-41d4-a716-446655440000",
					},
				},
				want:    &emptypb.Empty{},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServer{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
			}
			got, err := impl.DeletePost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServer.DeletePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServer.DeletePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guestbookServiceServerImpl_ListPosts(t *testing.T) {
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *domain.Config
		usecase                             Usecase
	}
	type args struct {
		ctx context.Context
		req *pb.ListPostsRequest
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *pb.ListPostsResponse
		wantErr bool
	}
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			config := &domain.Config{}
			posts := []*domain.Post{
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example1"), domain.NewPostBody(""), now, time.Time{}, time.Time{}),
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440001"), domain.NewPostBody("example2"), domain.NewPostBody(""), now, time.Time{}, time.Time{}),
			}
			mockUsecase := &UsecaseMock{
				RangeFunc: func(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, domain.PaginationID, error) {
					return posts, domain.PaginationID{}, nil
				},
			}
			return testCase{
				name: "normal without next page",
				fields: fields{
					config:  config,
					usecase: mockUsecase,
				},
				args: args{
					ctx: ctx,
					req: &pb.ListPostsRequest{
						PageSize: 10,
					},
				},
				want: &pb.ListPostsResponse{
					Posts: []*pb.Post{
						{PostId: "550e8400-e29b-41d4-a716-446655440000", Body: "example1", Valid: true, CreateTime: timestamppb.New(now)},
						{PostId: "550e8400-e29b-41d4-a716-446655440001", Body: "example2", Valid: true, CreateTime: timestamppb.New(now)},
					},
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			config := &domain.Config{}
			posts := []*domain.Post{
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example1"), domain.NewPostBody(""), now, time.Time{}, time.Time{}),
			}
			nextPaginationID := domain.NewPaginationID("660e8400-e29b-41d4-a716-446655440000")
			mockUsecase := &UsecaseMock{
				RangeFunc: func(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, domain.PaginationID, error) {
					return posts, nextPaginationID, nil
				},
			}
			return testCase{
				name: "normal with next page",
				fields: fields{
					config:  config,
					usecase: mockUsecase,
				},
				args: args{
					ctx: ctx,
					req: &pb.ListPostsRequest{
						PageSize: 1,
					},
				},
				want: &pb.ListPostsResponse{
					Posts: []*pb.Post{
						{PostId: "550e8400-e29b-41d4-a716-446655440000", Body: "example1", Valid: true, CreateTime: timestamppb.New(now)},
					},
					NextPageToken: "660e8400-e29b-41d4-a716-446655440000",
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			config := &domain.Config{}
			mockUsecase := &UsecaseMock{
				RangeFunc: func(ctx context.Context, pageOption *domain.PageOption) ([]*domain.Post, domain.PaginationID, error) {
					return nil, domain.PaginationID{}, domain.ErrInvalidArgument
				},
			}
			return testCase{
				name: "error",
				fields: fields{
					config:  config,
					usecase: mockUsecase,
				},
				args: args{
					ctx: ctx,
					req: &pb.ListPostsRequest{
						PageSize: 10,
					},
				},
				want:    nil,
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServer{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
			}
			got, err := impl.ListPosts(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServer.ListPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServer.ListPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}
