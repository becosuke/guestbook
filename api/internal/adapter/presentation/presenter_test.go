package presentation

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/becosuke/guestbook/api/internal/domain/entity/config"
	entity "github.com/becosuke/guestbook/api/internal/domain/entity/post"
	"github.com/becosuke/guestbook/api/internal/domain/repository"
	pb "github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func Test_guestbookServiceServerImpl_GetPost(t *testing.T) {
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *config.Config
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
			config := &config.Config{}
			postID := entity.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := entity.NewBody("example")
			post := entity.NewPost(postID, body)
			mockUsecase := &UsecaseMock{
				GetFunc: func(ctx context.Context, id *entity.PostID) (*entity.Post, error) {
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
					PostId: "550e8400-e29b-41d4-a716-446655440000",
					Body:   "example",
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			config := &config.Config{}
			mockUsecase := &UsecaseMock{
				GetFunc: func(ctx context.Context, id *entity.PostID) (*entity.Post, error) {
					return nil, repository.ErrNotFound
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
		config                              *config.Config
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
			config := &config.Config{}
			res := entity.NewPost(entity.NewPostID("550e8400-e29b-41d4-a716-446655440000"), entity.NewBody("example"))
			mockUsecase := &UsecaseMock{
				CreateFunc: func(ctx context.Context, post *entity.Post) (*entity.Post, error) {
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
					PostId: "550e8400-e29b-41d4-a716-446655440000",
					Body:   "example",
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
		config                              *config.Config
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
			config := &config.Config{}
			post := entity.NewPost(entity.NewPostID("550e8400-e29b-41d4-a716-446655440000"), entity.NewBody("example-value"))
			mockUsecase := &UsecaseMock{
				UpdateFunc: func(ctx context.Context, p *entity.Post) (*entity.Post, error) {
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
					PostId: "550e8400-e29b-41d4-a716-446655440000",
					Body:   "example-value",
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
		config                              *config.Config
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
			config := &config.Config{}
			mockUsecase := &UsecaseMock{
				DeleteFunc: func(ctx context.Context, id *entity.PostID) error {
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
