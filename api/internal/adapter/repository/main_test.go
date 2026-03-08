//go:build integration

package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/domain/interfaces"
)

var (
	testPool      *pgxpool.Pool
	testQuerier   interfaces.PostQuerier
	testCommander interfaces.PostCommander
	testPaginator interfaces.Paginator
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://guestbook:guestbook@localhost:5432/guestbook?sslmode=disable"
	}

	pool, err := NewPool(ctx, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	testPool = pool

	config := &domain.Config{
		Environment: domain.EnvTest,
		LogLevel:    zapcore.InfoLevel,
		DatabaseURL: databaseURL,
	}
	logger := zap.NewNop()

	testQuerier = NewPostQuerier(config, logger, pool)
	testCommander = NewPostCommander(config, logger, pool)
	testPaginator = NewPaginator(config, logger, pool)

	os.Exit(m.Run())
}

func truncateTables(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	_, err := testPool.Exec(ctx, `TRUNCATE TABLE Posts, Paginations`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func newUUID() string {
	return uuid.New().String()
}
