//go:generate moq -out post_commander_mock.go -pkg interfaces . PostCommander
package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

type PostCommander interface {
	Create(context.Context, *domain.Post) error
	Update(context.Context, *domain.Post) error
	Delete(context.Context, *domain.PostID) error
}
