//go:generate moq -out usecase_mock.go -pkg presentation . Usecase
package presentation

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

type Usecase interface {
	Get(context.Context, domain.PostID) (*domain.Post, error)
	Range(context.Context, *domain.PageOption) ([]*domain.Post, *domain.PaginationID, error)
	Create(context.Context, *domain.Post) (*domain.Post, error)
	Update(context.Context, *domain.Post) (*domain.Post, error)
	Delete(context.Context, domain.PostID) error
}
