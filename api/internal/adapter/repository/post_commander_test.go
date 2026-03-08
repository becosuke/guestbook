//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/becosuke/guestbook/api/internal/domain"
)

func TestCreatePost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		id := newUUID()
		post := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("hello world"), time.Time{}, nil)
		err := testCommander.CreatePost(ctx, post)
		require.NoError(t, err)

		got, err := testQuerier.GetPost(ctx, domain.NewPostID(id))
		require.NoError(t, err)
		assert.Equal(t, id, got.PostID().String())
		assert.Equal(t, "hello world", got.PostBody().String())
		assert.True(t, got.Valid())
	})

	t.Run("duplicate id", func(t *testing.T) {
		id := newUUID()
		post := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("first"), time.Time{}, nil)
		err := testCommander.CreatePost(ctx, post)
		require.NoError(t, err)

		dup := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("second"), time.Time{}, nil)
		err = testCommander.CreatePost(ctx, dup)
		assert.ErrorIs(t, err, domain.ErrAlreadyExists)
	})
}

func TestUpdatePost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		id := newUUID()
		post := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("original"), time.Time{}, nil)
		err := testCommander.CreatePost(ctx, post)
		require.NoError(t, err)

		updated := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("updated"), time.Time{}, nil)
		err = testCommander.UpdatePost(ctx, updated)
		require.NoError(t, err)

		got, err := testQuerier.GetPost(ctx, domain.NewPostID(id))
		require.NoError(t, err)
		assert.Equal(t, "updated", got.PostBody().String())
	})

	t.Run("not found", func(t *testing.T) {
		post := domain.NewPost(domain.NewPostID(newUUID()), domain.NewPostBody("body"), time.Time{}, nil)
		err := testCommander.UpdatePost(ctx, post)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("deleted post", func(t *testing.T) {
		id := newUUID()
		post := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("to delete"), time.Time{}, nil)
		err := testCommander.CreatePost(ctx, post)
		require.NoError(t, err)

		err = testCommander.DeletePost(ctx, domain.NewPostID(id))
		require.NoError(t, err)

		updated := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("should fail"), time.Time{}, nil)
		err = testCommander.UpdatePost(ctx, updated)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestDeletePost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		id := newUUID()
		post := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("to delete"), time.Time{}, nil)
		err := testCommander.CreatePost(ctx, post)
		require.NoError(t, err)

		err = testCommander.DeletePost(ctx, domain.NewPostID(id))
		require.NoError(t, err)

		got, err := testQuerier.GetPost(ctx, domain.NewPostID(id))
		require.NoError(t, err)
		assert.False(t, got.Valid())
		assert.NotNil(t, got.DeleteTime())
	})

	t.Run("not found", func(t *testing.T) {
		err := testCommander.DeletePost(ctx, domain.NewPostID(newUUID()))
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("double delete", func(t *testing.T) {
		id := newUUID()
		post := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("double delete"), time.Time{}, nil)
		err := testCommander.CreatePost(ctx, post)
		require.NoError(t, err)

		err = testCommander.DeletePost(ctx, domain.NewPostID(id))
		require.NoError(t, err)

		err = testCommander.DeletePost(ctx, domain.NewPostID(id))
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
