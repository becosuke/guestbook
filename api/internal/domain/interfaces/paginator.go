//go:generate moq -out paginator_mock.go -pkg interfaces . Paginator
package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

type Paginator interface {
	Get(context.Context, *domain.PaginationID) (*domain.Pagination, error)
	Save(context.Context, *domain.Pagination) error
}
