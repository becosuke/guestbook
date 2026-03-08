package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func truncateTables(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE Posts, Paginations")
	require.NoError(t, err, "failed to truncate tables")
}

func newUUID() string {
	return uuid.New().String()
}
