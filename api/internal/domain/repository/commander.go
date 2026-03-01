//go:generate moq -out commander_mock.go -pkg repository . Commander
package repository

import (
	"context"

	entity "github.com/becosuke/guestbook/api/internal/domain/entity/post"
)

type Commander interface {
	Create(context.Context, *entity.Post) error
	Update(context.Context, *entity.Post) error
	Delete(context.Context, *entity.PostID) error
}
