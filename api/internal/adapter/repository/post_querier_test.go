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

func TestGetPost(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("existing post", func(t *testing.T) {
		id := newUUID()
		body := "test body"
		post := domain.NewPost(domain.NewPostID(id), domain.NewPostBody(body), domain.NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
		err := testCommander.CreatePost(ctx, post)
		require.NoError(t, err)

		got, err := testQuerier.GetPost(ctx, domain.NewPostID(id))
		require.NoError(t, err)
		assert.Equal(t, id, got.PostID().String())
		assert.Equal(t, body, got.PostBody().String())
		assert.True(t, got.Valid())
		assert.False(t, got.CreateTime().IsZero())
	})

	t.Run("not found", func(t *testing.T) {
		got, err := testQuerier.GetPost(ctx, domain.NewPostID(newUUID()))
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, got)
	})

	t.Run("deleted post is still returned by GetPost", func(t *testing.T) {
		id := newUUID()
		post := domain.NewPost(domain.NewPostID(id), domain.NewPostBody("to be deleted"), domain.NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
		err := testCommander.CreatePost(ctx, post)
		require.NoError(t, err)

		err = testCommander.DeletePost(ctx, domain.NewPostID(id))
		require.NoError(t, err)

		got, err := testQuerier.GetPost(ctx, domain.NewPostID(id))
		require.NoError(t, err)
		assert.Equal(t, id, got.PostID().String())
		assert.False(t, got.Valid())
		assert.False(t, got.DeleteTime().IsZero())
	})
}

func TestRangePosts(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("empty result", func(t *testing.T) {
		posts, err := testQuerier.RangePosts(ctx, 10, nil)
		require.NoError(t, err)
		assert.Empty(t, posts)
	})

	t.Run("without cursor", func(t *testing.T) {
		ids := make([]string, 3)
		for i := range ids {
			ids[i] = newUUID()
			post := domain.NewPost(domain.NewPostID(ids[i]), domain.NewPostBody("body "+ids[i][:8]), domain.NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
			err := testCommander.CreatePost(ctx, post)
			require.NoError(t, err)
			// Sleep briefly to ensure distinct CreateTime values
			time.Sleep(10 * time.Millisecond)
		}

		posts, err := testQuerier.RangePosts(ctx, 10, nil)
		require.NoError(t, err)
		assert.Len(t, posts, 3)
		// Verify ordering: CreateTime DESC
		for i := 0; i < len(posts)-1; i++ {
			assert.True(t, posts[i].CreateTime().After(posts[i+1].CreateTime()) || posts[i].CreateTime().Equal(posts[i+1].CreateTime()))
		}
	})

	t.Run("with cursor", func(t *testing.T) {
		truncateTables(t)
		ids := make([]string, 5)
		for i := range ids {
			ids[i] = newUUID()
			post := domain.NewPost(domain.NewPostID(ids[i]), domain.NewPostBody("body "+ids[i][:8]), domain.NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
			err := testCommander.CreatePost(ctx, post)
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond)
		}

		// Get first page
		firstPage, err := testQuerier.RangePosts(ctx, 3, nil)
		require.NoError(t, err)
		require.Len(t, firstPage, 3)

		// Use last item as cursor for next page
		lastPost := firstPage[len(firstPage)-1]
		cursor := domain.NewPostCursor(lastPost.PostID().String(), lastPost.CreateTime())

		secondPage, err := testQuerier.RangePosts(ctx, 3, cursor)
		require.NoError(t, err)
		assert.Len(t, secondPage, 2)

		// Verify no overlap between pages
		firstPageIDs := make(map[string]bool)
		for _, p := range firstPage {
			firstPageIDs[p.PostID().String()] = true
		}
		for _, p := range secondPage {
			assert.False(t, firstPageIDs[p.PostID().String()], "pages should not overlap")
		}
	})

	t.Run("respects page size", func(t *testing.T) {
		truncateTables(t)
		for i := 0; i < 5; i++ {
			post := domain.NewPost(domain.NewPostID(newUUID()), domain.NewPostBody("body"), domain.NewPostBody(""), time.Time{}, time.Time{}, time.Time{})
			err := testCommander.CreatePost(ctx, post)
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond)
		}

		posts, err := testQuerier.RangePosts(ctx, 2, nil)
		require.NoError(t, err)
		assert.Len(t, posts, 2)
	})
}
