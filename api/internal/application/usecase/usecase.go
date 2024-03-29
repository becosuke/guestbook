package usecase

import (
	"context"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
)

type Usecase interface {
	Get(context.Context, *domain.Serial) (*domain.Post, error)
	Range(context.Context, *domain.PageOption) ([]*domain.Post, error)
	Create(context.Context, *domain.Post) (*domain.Post, error)
	Update(context.Context, *domain.Post) (*domain.Post, error)
	Delete(context.Context, *domain.Serial) error
}
