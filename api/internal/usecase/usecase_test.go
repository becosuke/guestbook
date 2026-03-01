package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain/entity/config"
	entity "github.com/becosuke/guestbook/api/internal/domain/entity/post"
	"github.com/becosuke/guestbook/api/internal/domain/repository"
)

func TestUsecase_Get(t *testing.T) {
	type fields struct {
		querier   repository.Querier
		commander repository.Commander
	}
	type args struct {
		ctx    context.Context
		postID *entity.PostID
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *entity.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			postID := entity.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := entity.NewBody("example")
			post := entity.NewPost(postID, body)
			mockQuerier := &repository.QuerierMock{
				GetFunc: func(ctx context.Context, id *entity.PostID) (*entity.Post, error) {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
					return post, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					querier: mockQuerier,
				},
				args: args{
					ctx:    ctx,
					postID: postID,
				},
				want:    post,
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			postID := entity.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			mockQuerier := &repository.QuerierMock{
				GetFunc: func(ctx context.Context, id *entity.PostID) (*entity.Post, error) {
					return nil, repository.ErrNotFound
				},
			}
			return testCase{
				name: "not found",
				fields: fields{
					querier: mockQuerier,
				},
				args: args{
					ctx:    ctx,
					postID: postID,
				},
				want:    nil,
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(&config.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
			got, err := uc.Get(tt.args.ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_Range(t *testing.T) {
	type fields struct {
		querier   repository.Querier
		commander repository.Commander
	}
	type args struct {
		ctx        context.Context
		pageOption *entity.PageOption
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    []*entity.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			pageSize := entity.PageSize(10)
			pageOption := entity.NewPageOption(&pageSize, nil)
			posts := []*entity.Post{
				entity.NewPost(entity.NewPostID("550e8400-e29b-41d4-a716-446655440000"), entity.NewBody("example1")),
				entity.NewPost(entity.NewPostID("550e8400-e29b-41d4-a716-446655440001"), entity.NewBody("example2")),
			}
			mockQuerier := &repository.QuerierMock{
				RangeFunc: func(ctx context.Context, po *entity.PageOption) ([]*entity.Post, error) {
					return posts, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					querier: mockQuerier,
				},
				args: args{
					ctx:        ctx,
					pageOption: pageOption,
				},
				want:    posts,
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			pageSize := entity.PageSize(10)
			pageOption := entity.NewPageOption(&pageSize, nil)
			mockQuerier := &repository.QuerierMock{
				RangeFunc: func(ctx context.Context, po *entity.PageOption) ([]*entity.Post, error) {
					return nil, repository.ErrInvalidArgument
				},
			}
			return testCase{
				name: "error",
				fields: fields{
					querier: mockQuerier,
				},
				args: args{
					ctx:        ctx,
					pageOption: pageOption,
				},
				want:    nil,
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(&config.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
			got, err := uc.Range(tt.args.ctx, tt.args.pageOption)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Range() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_Create(t *testing.T) {
	type fields struct {
		querier   repository.Querier
		commander repository.Commander
	}
	type args struct {
		ctx  context.Context
		post *entity.Post
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *entity.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			body := entity.NewBody("example")
			inputPost := entity.NewPost(nil, body)
			returnedPost := entity.NewPost(entity.NewPostID("550e8400-e29b-41d4-a716-446655440000"), body)
			mockCommander := &repository.CommanderMock{
				CreateFunc: func(ctx context.Context, p *entity.Post) error {
					assert.Equal(t, "example", p.Body().String())
					assert.NotEmpty(t, p.PostID().String())
					return nil
				},
			}
			mockQuerier := &repository.QuerierMock{
				GetFunc: func(ctx context.Context, id *entity.PostID) (*entity.Post, error) {
					return returnedPost, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					querier:   mockQuerier,
					commander: mockCommander,
				},
				args: args{
					ctx:  ctx,
					post: inputPost,
				},
				want:    returnedPost,
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			body := entity.NewBody("example")
			inputPost := entity.NewPost(nil, body)
			mockCommander := &repository.CommanderMock{
				CreateFunc: func(ctx context.Context, p *entity.Post) error {
					return repository.ErrAlreadyExists
				},
			}
			return testCase{
				name: "create error",
				fields: fields{
					commander: mockCommander,
				},
				args: args{
					ctx:  ctx,
					post: inputPost,
				},
				want:    nil,
				wantErr: true,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			body := entity.NewBody("example")
			inputPost := entity.NewPost(nil, body)
			mockCommander := &repository.CommanderMock{
				CreateFunc: func(ctx context.Context, p *entity.Post) error {
					return nil
				},
			}
			mockQuerier := &repository.QuerierMock{
				GetFunc: func(ctx context.Context, id *entity.PostID) (*entity.Post, error) {
					return nil, repository.ErrNotFound
				},
			}
			return testCase{
				name: "get error",
				fields: fields{
					querier:   mockQuerier,
					commander: mockCommander,
				},
				args: args{
					ctx:  ctx,
					post: inputPost,
				},
				want:    nil,
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(&config.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
			got, err := uc.Create(tt.args.ctx, tt.args.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_Update(t *testing.T) {
	type fields struct {
		querier   repository.Querier
		commander repository.Commander
	}
	type args struct {
		ctx  context.Context
		post *entity.Post
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *entity.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			postID := entity.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := entity.NewBody("updated-example")
			post := entity.NewPost(postID, body)
			mockCommander := &repository.CommanderMock{
				UpdateFunc: func(ctx context.Context, p *entity.Post) error {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", p.PostID().String())
					assert.Equal(t, "updated-example", p.Body().String())
					return nil
				},
			}
			mockQuerier := &repository.QuerierMock{
				GetFunc: func(ctx context.Context, id *entity.PostID) (*entity.Post, error) {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
					return post, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					querier:   mockQuerier,
					commander: mockCommander,
				},
				args: args{
					ctx:  ctx,
					post: post,
				},
				want:    post,
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			postID := entity.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := entity.NewBody("updated-example")
			post := entity.NewPost(postID, body)
			mockCommander := &repository.CommanderMock{
				UpdateFunc: func(ctx context.Context, p *entity.Post) error {
					return repository.ErrNotFound
				},
			}
			return testCase{
				name: "update error",
				fields: fields{
					commander: mockCommander,
				},
				args: args{
					ctx:  ctx,
					post: post,
				},
				want:    nil,
				wantErr: true,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			postID := entity.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := entity.NewBody("updated-example")
			post := entity.NewPost(postID, body)
			mockCommander := &repository.CommanderMock{
				UpdateFunc: func(ctx context.Context, p *entity.Post) error {
					return nil
				},
			}
			mockQuerier := &repository.QuerierMock{
				GetFunc: func(ctx context.Context, id *entity.PostID) (*entity.Post, error) {
					return nil, repository.ErrNotFound
				},
			}
			return testCase{
				name: "get error",
				fields: fields{
					querier:   mockQuerier,
					commander: mockCommander,
				},
				args: args{
					ctx:  ctx,
					post: post,
				},
				want:    nil,
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(&config.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
			got, err := uc.Update(tt.args.ctx, tt.args.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_Delete(t *testing.T) {
	type fields struct {
		querier   repository.Querier
		commander repository.Commander
	}
	type args struct {
		ctx    context.Context
		postID *entity.PostID
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			postID := entity.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			mockCommander := &repository.CommanderMock{
				DeleteFunc: func(ctx context.Context, id *entity.PostID) error {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
					return nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					commander: mockCommander,
				},
				args: args{
					ctx:    ctx,
					postID: postID,
				},
				wantErr: false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			postID := entity.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			mockCommander := &repository.CommanderMock{
				DeleteFunc: func(ctx context.Context, id *entity.PostID) error {
					return repository.ErrNotFound
				},
			}
			return testCase{
				name: "error",
				fields: fields{
					commander: mockCommander,
				},
				args: args{
					ctx:    ctx,
					postID: postID,
				},
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(&config.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
			err := uc.Delete(tt.args.ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
