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
	type args struct {
		ctx    context.Context
		postID domain.PostID
	}
	type testCase struct {
		name    string
		repos   interfaces.Repositories
		args    args
		want    *domain.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := domain.NewPostBody("example")
			post := domain.NewPost(postID, body, domain.NewPostBody(""), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}, time.Time{})
			return testCase{
				name: "normal",
				repos: &interfaces.RepositoriesMock{
					GetPostFunc: func(ctx context.Context, id domain.PostID) (*domain.Post, error) {
						assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
						return post, nil
					},
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
			return testCase{
				name: "not found",
				repos: &interfaces.RepositoriesMock{
					GetPostFunc: func(ctx context.Context, id domain.PostID) (*domain.Post, error) {
						return nil, domain.ErrNotFound
					},
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.repos)
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
	type args struct {
		ctx        context.Context
		pageOption *domain.PageOption
	}
	type testCase struct {
		name             string
		repos            interfaces.Repositories
		args             args
		want             []*domain.Post
		wantPaginationID bool
		wantErr          bool
	}
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			pageSize := domain.PageSize(10)
			pageOption := domain.NewPageOption(&pageSize, nil)
			posts := []*domain.Post{
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example1"), domain.NewPostBody(""), now, time.Time{}, time.Time{}),
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440001"), domain.NewPostBody("example2"), domain.NewPostBody(""), now, time.Time{}, time.Time{}),
			}
			return testCase{
				name: "no next page",
				repos: &interfaces.RepositoriesMock{
					RangePostsFunc: func(ctx context.Context, ps int32, cursor *domain.PostCursor) ([]*domain.Post, error) {
						assert.Equal(t, int32(11), ps)
						assert.Nil(t, cursor)
						return posts, nil
					},
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
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), domain.NewPostBody("example1"), domain.NewPostBody(""), now, time.Time{}, time.Time{}),
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440001"), domain.NewPostBody("example2"), domain.NewPostBody(""), now, time.Time{}, time.Time{}),
				domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440002"), domain.NewPostBody("example3"), domain.NewPostBody(""), now, time.Time{}, time.Time{}),
			}
			return testCase{
				name: "has next page",
				repos: &interfaces.RepositoriesMock{
					RangePostsFunc: func(ctx context.Context, ps int32, cursor *domain.PostCursor) ([]*domain.Post, error) {
						assert.Equal(t, int32(3), ps)
						return allPosts, nil
					},
					SavePaginationFunc: func(ctx context.Context, pagination *domain.Pagination) error {
						return nil
					},
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
			return testCase{
				name: "error",
				repos: &interfaces.RepositoriesMock{
					RangePostsFunc: func(ctx context.Context, ps int32, cursor *domain.PostCursor) ([]*domain.Post, error) {
						return nil, domain.ErrInvalidArgument
					},
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.repos)
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
	type args struct {
		ctx  context.Context
		post *domain.Post
	}
	type testCase struct {
		name    string
		repos   interfaces.Repositories
		args    args
		want    *domain.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			body := domain.NewPostBody("example")
			inputPost := domain.NewPost(domain.PostID{}, body, domain.NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
			returnedPost := domain.NewPost(domain.NewPostID("550e8400-e29b-41d4-a716-446655440000"), body, domain.NewPostBody(""), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}, time.Time{})
			return testCase{
				name: "normal",
				repos: &interfaces.RepositoriesMock{
					CreatePostFunc: func(ctx context.Context, p *domain.Post) error {
						assert.Equal(t, "example", p.PostBody().String())
						assert.NotEmpty(t, p.PostID().String())
						return nil
					},
					GetPostFunc: func(ctx context.Context, id domain.PostID) (*domain.Post, error) {
						return returnedPost, nil
					},
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
			inputPost := domain.NewPost(domain.PostID{}, body, domain.NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
			return testCase{
				name: "create error",
				repos: &interfaces.RepositoriesMock{
					CreatePostFunc: func(ctx context.Context, p *domain.Post) error {
						return domain.ErrAlreadyExists
					},
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
			inputPost := domain.NewPost(domain.PostID{}, body, domain.NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
			return testCase{
				name: "get error",
				repos: &interfaces.RepositoriesMock{
					CreatePostFunc: func(ctx context.Context, p *domain.Post) error {
						return nil
					},
					GetPostFunc: func(ctx context.Context, id domain.PostID) (*domain.Post, error) {
						return nil, domain.ErrNotFound
					},
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.repos)
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
	type args struct {
		ctx  context.Context
		post *domain.Post
	}
	type testCase struct {
		name    string
		repos   interfaces.Repositories
		args    args
		want    *domain.Post
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			body := domain.NewPostBody("updated-example")
			post := domain.NewPost(postID, body, domain.NewPostBody(""), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}, time.Time{})
			return testCase{
				name: "normal",
				repos: &interfaces.RepositoriesMock{
					UpdatePostFunc: func(ctx context.Context, p *domain.Post) error {
						assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", p.PostID().String())
						assert.Equal(t, "updated-example", p.PostBody().String())
						return nil
					},
					GetPostFunc: func(ctx context.Context, id domain.PostID) (*domain.Post, error) {
						assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
						return post, nil
					},
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
			post := domain.NewPost(postID, body, domain.NewPostBody(""), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}, time.Time{})
			return testCase{
				name: "update error",
				repos: &interfaces.RepositoriesMock{
					UpdatePostFunc: func(ctx context.Context, p *domain.Post) error {
						return domain.ErrNotFound
					},
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
			post := domain.NewPost(postID, body, domain.NewPostBody(""), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}, time.Time{})
			return testCase{
				name: "get error",
				repos: &interfaces.RepositoriesMock{
					UpdatePostFunc: func(ctx context.Context, p *domain.Post) error {
						return nil
					},
					GetPostFunc: func(ctx context.Context, id domain.PostID) (*domain.Post, error) {
						return nil, domain.ErrNotFound
					},
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.repos)
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
	type args struct {
		ctx    context.Context
		postID domain.PostID
	}
	type testCase struct {
		name    string
		repos   interfaces.Repositories
		args    args
		wantErr bool
	}
	tests := []testCase{
		func() testCase {
			ctx := context.Background()
			postID := domain.NewPostID("550e8400-e29b-41d4-a716-446655440000")
			return testCase{
				name: "normal",
				repos: &interfaces.RepositoriesMock{
					DeletePostFunc: func(ctx context.Context, id domain.PostID) error {
						assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
						return nil
					},
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
			return testCase{
				name: "error",
				repos: &interfaces.RepositoriesMock{
					DeletePostFunc: func(ctx context.Context, id domain.PostID) error {
						return domain.ErrNotFound
					},
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
			uc := NewUsecase(&domain.Config{}, zap.NewNop(), tt.repos)
			err := uc.Delete(tt.args.ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
