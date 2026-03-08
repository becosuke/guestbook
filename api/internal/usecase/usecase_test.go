package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

func TestUsecase_Get(t *testing.T) {
	type fields struct {
		postQuerier   interfaces.PostQuerier
		postCommander interfaces.PostCommander
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
			post := domain.NewPost(postID, body, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
			mockQuerier := &interfaces.PostQuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
					return post, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					postQuerier: mockQuerier,
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
			mockQuerier := &interfaces.PostQuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					return nil, domain.ErrNotFound
				},
			}
			return testCase{
				name: "not found",
				fields: fields{
					postQuerier: mockQuerier,
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.postQuerier, tt.fields.postCommander, nil)
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
		postQuerier interfaces.PostQuerier
		paginator   interfaces.Paginator
	}
	type args struct {
		ctx        context.Context
		pageOption *domain.PageOption
	}
	type testCase struct {
		name              string
		fields            fields
		args              args
		want              []*domain.Post
		wantPaginationID  bool
		wantErr           bool
	}
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			pageSize := domain.PageSize(10)
			pageOption := domain.NewPageOption(&pageSize, nil)
			posts := []*domain.Post{
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example1"), now, nil),
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440001"), domain.NewPostBody("example2"), now, nil),
			}
			mockQuerier := &interfaces.PostQuerierMock{
				RangeFunc: func(ctx context.Context, ps int32, cursor *domain.PostCursor) ([]*domain.Post, error) {
					assert.Equal(t, int32(11), ps)
					assert.Nil(t, cursor)
					return posts, nil
				},
			}
			return testCase{
				name: "no next page",
				fields: fields{
					postQuerier: mockQuerier,
				},
				args: args{
					ctx:        ctx,
					pageOption: pageOption,
				},
				want:             posts,
				wantPaginationID: false,
				wantErr:          false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			pageSize := domain.PageSize(2)
			pageOption := domain.NewPageOption(&pageSize, nil)
			allPosts := []*domain.Post{
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example1"), now, nil),
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440001"), domain.NewPostBody("example2"), now, nil),
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440002"), domain.NewPostBody("example3"), now, nil),
			}
			mockQuerier := &interfaces.PostQuerierMock{
				RangeFunc: func(ctx context.Context, ps int32, cursor *domain.PostCursor) ([]*domain.Post, error) {
					assert.Equal(t, int32(3), ps)
					return allPosts, nil
				},
			}
			mockPaginator := &interfaces.PaginatorMock{
				SaveFunc: func(ctx context.Context, pagination *domain.Pagination) error {
					return nil
				},
			}
			return testCase{
				name: "has next page",
				fields: fields{
					postQuerier: mockQuerier,
					paginator:   mockPaginator,
				},
				args: args{
					ctx:        ctx,
					pageOption: pageOption,
				},
				want:             allPosts[:2],
				wantPaginationID: true,
				wantErr:          false,
			}
		}(),
		func() testCase {
			ctx := context.Background()
			pageSize := domain.PageSize(10)
			pageOption := domain.NewPageOption(&pageSize, nil)
			mockQuerier := &interfaces.PostQuerierMock{
				RangeFunc: func(ctx context.Context, ps int32, cursor *domain.PostCursor) ([]*domain.Post, error) {
					return nil, domain.ErrInvalidArgument
				},
			}
			return testCase{
				name: "error",
				fields: fields{
					postQuerier: mockQuerier,
				},
				args: args{
					ctx:        ctx,
					pageOption: pageOption,
				},
				want:             nil,
				wantPaginationID: false,
				wantErr:          true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.postQuerier, nil, tt.fields.paginator)
			got, paginationID, err := uc.Range(tt.args.ctx, tt.args.pageOption)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Range() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Range() = %v, want %v", got, tt.want)
			}
			if tt.wantPaginationID {
				assert.NotNil(t, paginationID)
			} else {
				assert.Nil(t, paginationID)
			}
		})
	}
}

func TestUsecase_Create(t *testing.T) {
	type fields struct {
		postQuerier   interfaces.PostQuerier
		postCommander interfaces.PostCommander
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
			inputPost := domain.NewPost(nil, body, time.Time{}, nil)
			returnedPost := domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), body, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
			mockCommander := &interfaces.PostCommanderMock{
				CreateFunc: func(ctx context.Context, p *domain.Post) error {
					assert.Equal(t, "example", p.PostBody().String())
					assert.NotEmpty(t, p.PostID().String())
					return nil
				},
			}
			mockQuerier := &interfaces.PostQuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					return returnedPost, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					postQuerier:   mockQuerier,
					postCommander: mockCommander,
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
			inputPost := domain.NewPost(nil, body, time.Time{}, nil)
			mockCommander := &interfaces.PostCommanderMock{
				CreateFunc: func(ctx context.Context, p *domain.Post) error {
					return domain.ErrAlreadyExists
				},
			}
			return testCase{
				name: "create error",
				fields: fields{
					postCommander: mockCommander,
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
			inputPost := domain.NewPost(nil, body, time.Time{}, nil)
			mockCommander := &interfaces.PostCommanderMock{
				CreateFunc: func(ctx context.Context, p *domain.Post) error {
					return nil
				},
			}
			mockQuerier := &interfaces.PostQuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					return nil, domain.ErrNotFound
				},
			}
			return testCase{
				name: "get error",
				fields: fields{
					postQuerier:   mockQuerier,
					postCommander: mockCommander,
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.postQuerier, tt.fields.postCommander, nil)
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
		postQuerier   interfaces.PostQuerier
		postCommander interfaces.PostCommander
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
			post := domain.NewPost(postID, body, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
			mockCommander := &interfaces.PostCommanderMock{
				UpdateFunc: func(ctx context.Context, p *domain.Post) error {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", p.PostID().String())
					assert.Equal(t, "updated-example", p.PostBody().String())
					return nil
				},
			}
			mockQuerier := &interfaces.PostQuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
					return post, nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					postQuerier:   mockQuerier,
					postCommander: mockCommander,
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
			post := domain.NewPost(postID, body, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
			mockCommander := &interfaces.PostCommanderMock{
				UpdateFunc: func(ctx context.Context, p *domain.Post) error {
					return domain.ErrNotFound
				},
			}
			return testCase{
				name: "update error",
				fields: fields{
					postCommander: mockCommander,
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
			post := domain.NewPost(postID, body, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
			mockCommander := &interfaces.PostCommanderMock{
				UpdateFunc: func(ctx context.Context, p *domain.Post) error {
					return nil
				},
			}
			mockQuerier := &interfaces.PostQuerierMock{
				GetFunc: func(ctx context.Context, id *domain.PostID) (*domain.Post, error) {
					return nil, domain.ErrNotFound
				},
			}
			return testCase{
				name: "get error",
				fields: fields{
					postQuerier:   mockQuerier,
					postCommander: mockCommander,
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.fields.postQuerier, tt.fields.postCommander, nil)
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
		postCommander interfaces.PostCommander
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
			mockCommander := &interfaces.PostCommanderMock{
				DeleteFunc: func(ctx context.Context, id *domain.PostID) error {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
					return nil
				},
			}
			return testCase{
				name: "normal",
				fields: fields{
					postCommander: mockCommander,
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
			mockCommander := &interfaces.PostCommanderMock{
				DeleteFunc: func(ctx context.Context, id *domain.PostID) error {
					return domain.ErrNotFound
				},
			}
			return testCase{
				name: "error",
				fields: fields{
					postCommander: mockCommander,
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), nil, tt.fields.postCommander, nil)
			err := uc.Delete(tt.args.ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
