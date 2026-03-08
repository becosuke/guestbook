//go:generate moq -out post_commander_mock.go -pkg interfaces . PostCommander
package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

type PostCommander interface {
	CreatePost(context.Context, *domain.Post) error
	UpdatePost(context.Context, *domain.Post) error
	DeletePost(context.Context, *domain.PostID) error
}
