//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/becosuke/guestbook/api/internal/domain"
)

func TestGetPagination(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("existing pagination", func(t *testing.T) {
		id := newUUID()
		cursor := []byte(`{"last_post_id":"abc","last_create_time":"2024-01-01T00:00:00Z"}`)
		pagination := domain.NewPagination(domain.NewPaginationID(id), cursor)
		err := testPaginator.SavePagination(ctx, pagination)
		require.NoError(t, err)

		got, err := testPaginator.GetPagination(ctx, domain.NewPaginationID(id))
		require.NoError(t, err)
		assert.Equal(t, id, got.PaginationID().String())
		assert.Equal(t, cursor, got.Cursor())
	})

	t.Run("not found", func(t *testing.T) {
		got, err := testPaginator.GetPagination(ctx, domain.NewPaginationID(newUUID()))
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, got)
	})
}

func TestSavePagination(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		id := newUUID()
		cursor := []byte(`{"last_post_id":"xyz","last_create_time":"2024-06-15T12:00:00Z"}`)
		pagination := domain.NewPagination(domain.NewPaginationID(id), cursor)

		err := testPaginator.SavePagination(ctx, pagination)
		require.NoError(t, err)

		got, err := testPaginator.GetPagination(ctx, domain.NewPaginationID(id))
		require.NoError(t, err)
		assert.Equal(t, id, got.PaginationID().String())
		assert.Equal(t, cursor, got.Cursor())
	})
}
