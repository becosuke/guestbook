package test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func TestCreatePost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	body := "Hello, Guestbook!"
	resp, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
		Post:           &pb.Post{PostId: uuid.Nil.String(), Body: body},
		IdempotencyKey: newUUID(),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetPostId())

	_, parseErr := uuid.Parse(resp.GetPostId())
	assert.NoError(t, parseErr, "post_id should be a valid UUID")
	assert.Equal(t, body, resp.GetBody())
	assert.True(t, resp.GetValid())
}

func TestGetPost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("existing post", func(t *testing.T) {
		created, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			Post:           &pb.Post{PostId: uuid.Nil.String(), Body: "get me"},
			IdempotencyKey: newUUID(),
		})
		require.NoError(t, err)

		got, err := testClient.GetPost(ctx, &pb.GetPostRequest{
			PostId: created.GetPostId(),
		})
		require.NoError(t, err)
		assert.Equal(t, created.GetPostId(), got.GetPostId())
		assert.Equal(t, "get me", got.GetBody())
		assert.True(t, got.GetValid())
	})

	t.Run("not found", func(t *testing.T) {
		_, err := testClient.GetPost(ctx, &pb.GetPostRequest{
			PostId: newUUID(),
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})
}

func TestUpdatePost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	created, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
		Post:           &pb.Post{PostId: uuid.Nil.String(), Body: "original body"},
		IdempotencyKey: newUUID(),
	})
	require.NoError(t, err)

	t.Run("update body", func(t *testing.T) {
		updated, err := testClient.UpdatePost(ctx, &pb.UpdatePostRequest{
			Post: &pb.Post{
				PostId: created.GetPostId(),
				Body:   "updated body",
			},
			IdempotencyKey: newUUID(),
		})
		require.NoError(t, err)
		assert.Equal(t, created.GetPostId(), updated.GetPostId())
		assert.Equal(t, "updated body", updated.GetBody())
		assert.True(t, updated.GetValid())
	})

	t.Run("get after update returns new value", func(t *testing.T) {
		got, err := testClient.GetPost(ctx, &pb.GetPostRequest{
			PostId: created.GetPostId(),
		})
		require.NoError(t, err)
		assert.Equal(t, "updated body", got.GetBody())
	})
}

func TestDeletePost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	created, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
		Post:           &pb.Post{PostId: uuid.Nil.String(), Body: "to be deleted"},
		IdempotencyKey: newUUID(),
	})
	require.NoError(t, err)

	t.Run("delete post", func(t *testing.T) {
		_, err := testClient.DeletePost(ctx, &pb.DeletePostRequest{
			PostId:         created.GetPostId(),
			IdempotencyKey: newUUID(),
		})
		require.NoError(t, err)
	})

	t.Run("get after delete returns valid=false", func(t *testing.T) {
		got, err := testClient.GetPost(ctx, &pb.GetPostRequest{
			PostId: created.GetPostId(),
		})
		require.NoError(t, err)
		assert.False(t, got.GetValid())
	})

	t.Run("re-delete returns error", func(t *testing.T) {
		_, err := testClient.DeletePost(ctx, &pb.DeletePostRequest{
			PostId:         created.GetPostId(),
			IdempotencyKey: newUUID(),
		})
		require.Error(t, err)
	})
}

func TestListPosts(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("empty list", func(t *testing.T) {
		resp, err := testClient.ListPosts(ctx, &pb.ListPostsRequest{
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Empty(t, resp.GetPosts())
		assert.Empty(t, resp.GetNextPageToken())
	})

	t.Run("list with pagination", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			_, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
				Post:           &pb.Post{PostId: uuid.Nil.String(), Body: "post for list"},
				IdempotencyKey: newUUID(),
			})
			require.NoError(t, err)
		}

		allResp, err := testClient.ListPosts(ctx, &pb.ListPostsRequest{
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Len(t, allResp.GetPosts(), 5)
		assert.Empty(t, allResp.GetNextPageToken())

		seen := make(map[string]bool)
		pageToken := ""
		totalCollected := 0

		for {
			resp, err := testClient.ListPosts(ctx, &pb.ListPostsRequest{
				PageSize:  2,
				PageToken: pageToken,
			})
			require.NoError(t, err)

			for _, p := range resp.GetPosts() {
				assert.False(t, seen[p.GetPostId()], "duplicate post_id: %s", p.GetPostId())
				seen[p.GetPostId()] = true
			}
			totalCollected += len(resp.GetPosts())

			if resp.GetNextPageToken() == "" {
				break
			}
			pageToken = resp.GetNextPageToken()
		}

		assert.Equal(t, 5, totalCollected)
	})
}

func TestCRUDLifecycle(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	created, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
		Post:           &pb.Post{PostId: uuid.Nil.String(), Body: "lifecycle test"},
		IdempotencyKey: newUUID(),
	})
	require.NoError(t, err)
	assert.True(t, created.GetValid())

	got, err := testClient.GetPost(ctx, &pb.GetPostRequest{
		PostId: created.GetPostId(),
	})
	require.NoError(t, err)
	assert.Equal(t, "lifecycle test", got.GetBody())

	updated, err := testClient.UpdatePost(ctx, &pb.UpdatePostRequest{
		Post: &pb.Post{
			PostId: created.GetPostId(),
			Body:   "lifecycle updated",
		},
		IdempotencyKey: newUUID(),
	})
	require.NoError(t, err)
	assert.Equal(t, "lifecycle updated", updated.GetBody())

	gotAfterUpdate, err := testClient.GetPost(ctx, &pb.GetPostRequest{
		PostId: created.GetPostId(),
	})
	require.NoError(t, err)
	assert.Equal(t, "lifecycle updated", gotAfterUpdate.GetBody())

	listResp, err := testClient.ListPosts(ctx, &pb.ListPostsRequest{
		PageSize: 10,
	})
	require.NoError(t, err)
	found := false
	for _, p := range listResp.GetPosts() {
		if p.GetPostId() == created.GetPostId() {
			found = true
			break
		}
	}
	assert.True(t, found, "created post should appear in list")

	_, err = testClient.DeletePost(ctx, &pb.DeletePostRequest{
		PostId:         created.GetPostId(),
		IdempotencyKey: newUUID(),
	})
	require.NoError(t, err)

	gotAfterDelete, err := testClient.GetPost(ctx, &pb.GetPostRequest{
		PostId: created.GetPostId(),
	})
	require.NoError(t, err)
	assert.False(t, gotAfterDelete.GetValid())
}

// requireInvalidArgument asserts that err is a gRPC InvalidArgument status
// and returns the BadRequest details for further field violation assertions.
func requireInvalidArgument(t *testing.T, err error) *errdetails.BadRequest {
	t.Helper()
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	for _, d := range st.Details() {
		if br, ok := d.(*errdetails.BadRequest); ok {
			return br
		}
	}
	t.Fatal("expected BadRequest details in error")
	return nil
}

func TestGetPost_Validation(t *testing.T) {
	ctx := context.Background()

	t.Run("empty post_id", func(t *testing.T) {
		_, err := testClient.GetPost(ctx, &pb.GetPostRequest{PostId: ""})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("invalid UUID post_id", func(t *testing.T) {
		_, err := testClient.GetPost(ctx, &pb.GetPostRequest{PostId: "not-a-uuid"})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})
}

func TestCreatePost_Validation(t *testing.T) {
	ctx := context.Background()

	t.Run("missing post", func(t *testing.T) {
		_, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("empty body", func(t *testing.T) {
		_, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			Post:           &pb.Post{Body: ""},
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("body exceeds max length", func(t *testing.T) {
		longBody := strings.Repeat("a", 129)
		_, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			Post:           &pb.Post{Body: longBody},
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("body at max length is valid", func(t *testing.T) {
		truncateTables(t)
		maxBody := strings.Repeat("a", 128)
		resp, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			Post:           &pb.Post{PostId: uuid.Nil.String(), Body: maxBody},
			IdempotencyKey: newUUID(),
		})
		require.NoError(t, err)
		assert.Equal(t, maxBody, resp.GetBody())
	})

	t.Run("non-nil UUID post_id", func(t *testing.T) {
		_, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			Post:           &pb.Post{PostId: newUUID(), Body: "hello"},
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("invalid idempotency_key", func(t *testing.T) {
		_, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			Post:           &pb.Post{Body: "hello"},
			IdempotencyKey: "not-a-uuid",
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("empty idempotency_key", func(t *testing.T) {
		_, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			Post:           &pb.Post{Body: "hello"},
			IdempotencyKey: "",
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("multiple violations", func(t *testing.T) {
		_, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
			Post:           &pb.Post{Body: ""},
			IdempotencyKey: "not-a-uuid",
		})
		br := requireInvalidArgument(t, err)
		assert.GreaterOrEqual(t, len(br.GetFieldViolations()), 2,
			"should report violations for both body and idempotency_key")
	})
}

func TestUpdatePost_Validation(t *testing.T) {
	ctx := context.Background()

	t.Run("missing post", func(t *testing.T) {
		_, err := testClient.UpdatePost(ctx, &pb.UpdatePostRequest{
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("empty post_id", func(t *testing.T) {
		_, err := testClient.UpdatePost(ctx, &pb.UpdatePostRequest{
			Post: &pb.Post{
				PostId: "",
				Body:   "valid body",
			},
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("empty body", func(t *testing.T) {
		_, err := testClient.UpdatePost(ctx, &pb.UpdatePostRequest{
			Post: &pb.Post{
				PostId: newUUID(),
				Body:   "",
			},
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("body exceeds max length", func(t *testing.T) {
		_, err := testClient.UpdatePost(ctx, &pb.UpdatePostRequest{
			Post: &pb.Post{
				PostId: newUUID(),
				Body:   strings.Repeat("b", 129),
			},
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("invalid idempotency_key", func(t *testing.T) {
		_, err := testClient.UpdatePost(ctx, &pb.UpdatePostRequest{
			Post: &pb.Post{
				PostId: newUUID(),
				Body:   "valid body",
			},
			IdempotencyKey: "bad",
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})
}

func TestDeletePost_Validation(t *testing.T) {
	ctx := context.Background()

	t.Run("empty post_id", func(t *testing.T) {
		_, err := testClient.DeletePost(ctx, &pb.DeletePostRequest{
			PostId:         "",
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("invalid UUID post_id", func(t *testing.T) {
		_, err := testClient.DeletePost(ctx, &pb.DeletePostRequest{
			PostId:         "not-a-uuid",
			IdempotencyKey: newUUID(),
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("invalid idempotency_key", func(t *testing.T) {
		_, err := testClient.DeletePost(ctx, &pb.DeletePostRequest{
			PostId:         newUUID(),
			IdempotencyKey: "bad",
		})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("both fields invalid", func(t *testing.T) {
		_, err := testClient.DeletePost(ctx, &pb.DeletePostRequest{
			PostId:         "xxx",
			IdempotencyKey: "yyy",
		})
		br := requireInvalidArgument(t, err)
		assert.GreaterOrEqual(t, len(br.GetFieldViolations()), 2)
	})
}

func TestListPosts_Validation(t *testing.T) {
	ctx := context.Background()

	t.Run("zero page_size", func(t *testing.T) {
		_, err := testClient.ListPosts(ctx, &pb.ListPostsRequest{PageSize: 0})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})

	t.Run("negative page_size", func(t *testing.T) {
		_, err := testClient.ListPosts(ctx, &pb.ListPostsRequest{PageSize: -5})
		br := requireInvalidArgument(t, err)
		assert.NotEmpty(t, br.GetFieldViolations())
	})
}
