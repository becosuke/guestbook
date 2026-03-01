//go:generate moq -out commander_mock.go -pkg interfaces . Commander
package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain/entity"
)

type Commander interface {
	Create(context.Context, *entity.Post) error
	Update(context.Context, *entity.Post) error
	Delete(context.Context, *entity.PostID) error
}
