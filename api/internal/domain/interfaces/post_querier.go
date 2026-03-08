//go:generate moq -out post_querier_mock.go -pkg interfaces . PostQuerier
package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

type PostQuerier interface {
	GetPost(context.Context, *domain.PostID) (*domain.Post, error)
	RangePosts(ctx context.Context, pageSize int32, cursor *domain.PostCursor) ([]*domain.Post, error)
}
