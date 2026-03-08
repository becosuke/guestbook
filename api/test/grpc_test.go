package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

func TestCreatePost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	body := "Hello, Guestbook!"
	resp, err := testClient.CreatePost(ctx, &pb.CreatePostRequest{
		Post:           &pb.Post{Body: body},
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
			Post:           &pb.Post{Body: "get me"},
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
		Post:           &pb.Post{Body: "original body"},
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
		Post:           &pb.Post{Body: "to be deleted"},
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
				Post:           &pb.Post{Body: "post for list"},
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
		Post:           &pb.Post{Body: "lifecycle test"},
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
