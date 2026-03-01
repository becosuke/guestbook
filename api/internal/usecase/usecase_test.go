package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

func TestUsecase_Get(t *testing.T) {
	type fields struct {
		querier   interfaces.Querier
		commander interfaces.Commander
	}
	type args struct {
		ctx    context.Context
		postID *domain.PostID
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *domain.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := domain.NewPostBody("example")
			post := domain.NewPost(postID, body)
			mockQuerier := &interfaces.QuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
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
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			mockQuerier := &interfaces.QuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					return nil, domain.ErrNotFound
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
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
		querier   interfaces.Querier
		commander interfaces.Commander
	}
	type args struct {
		ctx        context.Context
		pageOption *domain.PageOption
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    []*domain.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			pageSize := domain.PageSize(10)
			pageOption := domain.NewPageOption(&pageSize, nil)
			posts := []*domain.Post{
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example1")),
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440001"), domain.NewPostBody("example2")),
			}
			mockQuerier := &interfaces.QuerierMock{
				RangeFunc: func(ctx context.Context, po *domain.PageOption) ([]*domain.Post, error) {
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
			pageSize := domain.PageSize(10)
			pageOption := domain.NewPageOption(&pageSize, nil)
			mockQuerier := &interfaces.QuerierMock{
				RangeFunc: func(ctx context.Context, po *domain.PageOption) ([]*domain.Post, error) {
					return nil, domain.ErrInvalidArgument
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
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
		querier   interfaces.Querier
		commander interfaces.Commander
	}
	type args struct {
		ctx  context.Context
		post *domain.Post
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *domain.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			body := domain.NewPostBody("example")
			inputPost := domain.NewPost(nil, body)
			returnedPost := domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), body)
			mockCommander := &interfaces.CommanderMock{
				CreateFunc: func(ctx context.Context, p *domain.Post) error {
					assert.Equal(t, "example", p.PostBody().String())
					assert.NotEmpty(t, p.PostID().String())
					return nil
				},
			}
			mockQuerier := &interfaces.QuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
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
			body := domain.NewPostBody("example")
			inputPost := domain.NewPost(nil, body)
			mockCommander := &interfaces.CommanderMock{
				CreateFunc: func(ctx context.Context, p *domain.Post) error {
					return domain.ErrAlreadyExists
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
			body := domain.NewPostBody("example")
			inputPost := domain.NewPost(nil, body)
			mockCommander := &interfaces.CommanderMock{
				CreateFunc: func(ctx context.Context, p *domain.Post) error {
					return nil
				},
			}
			mockQuerier := &interfaces.QuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					return nil, domain.ErrNotFound
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
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
		querier   interfaces.Querier
		commander interfaces.Commander
	}
	type args struct {
		ctx  context.Context
		post *domain.Post
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    *domain.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := domain.NewPostBody("updated-example")
			post := domain.NewPost(postID, body)
			mockCommander := &interfaces.CommanderMock{
				UpdateFunc: func(ctx context.Context, p *domain.Post) error {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", p.PostID().String())
					assert.Equal(t, "updated-example", p.PostBody().String())
					return nil
				},
			}
			mockQuerier := &interfaces.QuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
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
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := domain.NewPostBody("updated-example")
			post := domain.NewPost(postID, body)
			mockCommander := &interfaces.CommanderMock{
				UpdateFunc: func(ctx context.Context, p *domain.Post) error {
					return domain.ErrNotFound
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
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := domain.NewPostBody("updated-example")
			post := domain.NewPost(postID, body)
			mockCommander := &interfaces.CommanderMock{
				UpdateFunc: func(ctx context.Context, p *domain.Post) error {
					return nil
				},
			}
			mockQuerier := &interfaces.QuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					return nil, domain.ErrNotFound
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
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
		querier   interfaces.Querier
		commander interfaces.Commander
	}
	type args struct {
		ctx    context.Context
		postID *domain.PostID
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
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			mockCommander := &interfaces.CommanderMock{
				DeleteFunc: func(ctx context.Context, id *domain.PostID) error {
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
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			mockCommander := &interfaces.CommanderMock{
				DeleteFunc: func(ctx context.Context, id *domain.PostID) error {
					return domain.ErrNotFound
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.querier, tt.fields.commander)
			err := uc.Delete(tt.args.ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
