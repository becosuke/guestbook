package controller

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	repository "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	real_usecase "github.com/becosuke/guestbook/api/internal/application/usecase"
	domain "github.com/becosuke/guestbook/api/internal/domain/post"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
	pb "github.com/becosuke/guestbook/api/internal/pkg/pb"
	mock_usecase "github.com/becosuke/guestbook/api/mock/application/usecase"
)

func Test_guestbookServiceServerImpl_GetPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *pkgconfig.Config
		usecase                             real_usecase.Usecase
		boundary                            Boundary
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
			config := pkgconfig.NewConfig(ctx)
			boundary := NewBoundary()
			mockUsecase := mock_usecase.NewMockUsecase(ctrl)
			serial := domain.NewSerial(100)
			body := domain.NewBody("example")
			post := domain.NewPost(serial, body)
			mockUsecase.EXPECT().Get(ctx, serial).
				Return(post, nil).
				Do(func(ctx context.Context, serial *domain.Serial) {
					assert.Equal(t, int64(100), serial.Int64())
				})
			return testCase{
				name: "normal",
				fields: fields{
					config:   config,
					usecase:  mockUsecase,
					boundary: boundary,
				},
				args: args{
					ctx: ctx,
					req: &pb.GetPostRequest{
						Serial: 100,
					},
				},
				want: &pb.Post{
					Serial: 100,
					Body:   "example",
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			config := pkgconfig.NewConfig(ctx)
			boundary := NewBoundary()
			mockUsecase := mock_usecase.NewMockUsecase(ctrl)
			serial := domain.NewSerial(100)
			mockUsecase.EXPECT().Get(ctx, serial).Return(nil, repository.ErrMessageNotFound)
			return testCase{
				name: "not found",
				fields: fields{
					config:   config,
					usecase:  mockUsecase,
					boundary: boundary,
				},
				args: args{
					ctx: ctx,
					req: &pb.GetPostRequest{
						Serial: 100,
					},
				},
				want:    nil,
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServerImpl{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
				boundary:                            tt.fields.boundary,
			}
			got, err := impl.GetPost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServerImpl.GetPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServerImpl.GetPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guestbookServiceServerImpl_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *pkgconfig.Config
		usecase                             real_usecase.Usecase
		boundary                            Boundary
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
			config := pkgconfig.NewConfig(ctx)
			boundary := NewBoundary()
			mockUsecase := mock_usecase.NewMockUsecase(ctrl)
			req := domain.NewPost(domain.NewSerial(0), domain.NewBody("example"))
			res := domain.NewPost(domain.NewSerial(1), domain.NewBody("example"))
			mockUsecase.EXPECT().Create(ctx, req).Return(res, nil)
			return testCase{
				name: "normal",
				fields: fields{
					config:   config,
					usecase:  mockUsecase,
					boundary: boundary,
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
					Serial: 1,
					Body:   "example",
				},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServerImpl{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
				boundary:                            tt.fields.boundary,
			}
			got, err := impl.CreatePost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServerImpl.CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServerImpl.CreatePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guestbookServiceServerImpl_UpdatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *pkgconfig.Config
		usecase                             real_usecase.Usecase
		boundary                            Boundary
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
			config := pkgconfig.NewConfig(ctx)
			boundary := NewBoundary()
			mockUsecase := mock_usecase.NewMockUsecase(ctrl)
			post := domain.NewPost(domain.NewSerial(100), domain.NewBody("example-value"))
			mockUsecase.EXPECT().Update(ctx, post).Return(post, nil)
			return testCase{
				name: "normal",
				fields: fields{
					config:   config,
					usecase:  mockUsecase,
					boundary: boundary,
				},
				args: args{
					ctx: ctx,
					req: &pb.UpdatePostRequest{
						Post: &pb.Post{
							Serial: 100,
							Body:   "example-value",
						},
					},
				},
				want: &pb.Post{
					Serial: 100,
					Body:   "example-value",
				},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServerImpl{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
				boundary:                            tt.fields.boundary,
			}
			got, err := impl.UpdatePost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServerImpl.UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServerImpl.UpdatePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guestbookServiceServerImpl_DeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type fields struct {
		UnimplementedGuestbookServiceServer pb.UnimplementedGuestbookServiceServer
		config                              *pkgconfig.Config
		usecase                             real_usecase.Usecase
		boundary                            Boundary
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
			config := pkgconfig.NewConfig(ctx)
			boundary := NewBoundary()
			mockUsecase := mock_usecase.NewMockUsecase(ctrl)
			serial := domain.NewSerial(100)
			mockUsecase.EXPECT().Delete(ctx, serial).Return(nil)
			return testCase{
				name: "normal",
				fields: fields{
					config:   config,
					usecase:  mockUsecase,
					boundary: boundary,
				},
				args: args{
					ctx: ctx,
					req: &pb.DeletePostRequest{
						Serial: 100,
					},
				},
				want:    &emptypb.Empty{},
				wantErr: false,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &guestbookServiceServerImpl{
				UnimplementedGuestbookServiceServer: tt.fields.UnimplementedGuestbookServiceServer,
				config:                              tt.fields.config,
				usecase:                             tt.fields.usecase,
				boundary:                            tt.fields.boundary,
			}
			got, err := impl.DeletePost(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("guestbookServiceServerImpl.DeletePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guestbookServiceServerImpl.DeletePost() = %v, want %v", got, tt.want)
			}
		})
	}
}
