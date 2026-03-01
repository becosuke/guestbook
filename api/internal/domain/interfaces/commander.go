//go:generate moq -out commander_mock.go -pkg interfaces . Commander
package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

type Commander interface {
	Create(context.Context, *domain.Post) error
	Update(context.Context, *domain.Post) error
	Delete(context.Context, *domain.PostID) error
}
