package controller

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/becosuke/guestbook/api/internal/application/usecase"
	"github.com/becosuke/guestbook/api/internal/domain/post"
	"github.com/becosuke/guestbook/api/internal/registry/config"
	mock_usecase "github.com/becosuke/guestbook/api/mock/application/usecase"
	"github.com/becosuke/guestbook/pbgo"
)

func Test_guestbookServiceServerImpl_GetPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := config.NewConfig()
	mockBoundary := NewBoundary()
	type fields struct {
		UnimplementedGuestbookServiceServer pbgo.UnimplementedGuestbookServiceServer
		config                              *config.Config
		usecase                             usecase.Usecase
		boundary                            Boundary
	}
	type args struct {
		ctx context.Context
		req *pbgo.GetPostRequest
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *pbgo.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			mockUsecase := mock_usecase.NewMockUsecase(ctrl)
			serial := post.NewSerial(1)
			body := post.NewBody("example")
			mockUsecase.EXPECT().Get(ctx, serial).
				Return(post.NewPost(serial, body), nil).
				Do(func(ctx context.Context, serial *post.Serial) {
					assert.Equal(t, int64(1), serial.Int64())
				})
			return testCase{
				name: "normal",
				fields: fields{
					config:   mockConfig,
					usecase:  mockUsecase,
					boundary: mockBoundary,
				},
				args: args{
					ctx: ctx,
					req: &pbgo.GetPostRequest{
						Serial: 1,
					},
				},
				want: &pbgo.Post{
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
	type fields struct {
		UnimplementedGuestbookServiceServer pbgo.UnimplementedGuestbookServiceServer
		config                              *config.Config
		usecase                             usecase.Usecase
		boundary                            Boundary
	}
	type args struct {
		ctx context.Context
		req *pbgo.CreatePostRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pbgo.Post
		wantErr bool
	}{
		// TODO: Add test cases.
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
	type fields struct {
		UnimplementedGuestbookServiceServer pbgo.UnimplementedGuestbookServiceServer
		config                              *config.Config
		usecase                             usecase.Usecase
		boundary                            Boundary
	}
	type args struct {
		ctx context.Context
		req *pbgo.UpdatePostRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pbgo.Post
		wantErr bool
	}{
		// TODO: Add test cases.
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
	type fields struct {
		UnimplementedGuestbookServiceServer pbgo.UnimplementedGuestbookServiceServer
		config                              *config.Config
		usecase                             usecase.Usecase
		boundary                            Boundary
	}
	type args struct {
		ctx context.Context
		req *pbgo.DeletePostRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
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
