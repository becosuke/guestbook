package repository

import (
	"context"

	domain "github.com/becosuke/guestbook/api/internal/domain/post"
)

type Commander interface {
	Create(context.Context, *domain.Post) (*domain.Serial, error)
	Update(context.Context, *domain.Post) error
	Delete(context.Context, *domain.Serial) error
}
