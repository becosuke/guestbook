package repository

import (
	"context"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
)

type Querier interface {
	Get(context.Context, *domain.Serial) (*domain.Post, error)
	Range(context.Context, *domain.PageOption) ([]*domain.Post, error)
}
