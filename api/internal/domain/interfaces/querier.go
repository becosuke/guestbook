//go:generate moq -out querier_mock.go -pkg interfaces . Querier
package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

type Querier interface {
	Get(context.Context, *domain.PostID) (*domain.Post, error)
	Range(context.Context, *domain.PageOption) ([]*domain.Post, error)
}
