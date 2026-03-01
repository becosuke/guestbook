//go:generate moq -out usecase_mock.go -pkg presentation . Usecase
package presentation

import (
	"context"

	entity "github.com/becosuke/guestbook/api/internal/domain/entity/post"
)

type Usecase interface {
	Get(context.Context, *entity.PostID) (*entity.Post, error)
	Range(context.Context, *entity.PageOption) ([]*entity.Post, error)
	Create(context.Context, *entity.Post) (*entity.Post, error)
	Update(context.Context, *entity.Post) (*entity.Post, error)
	Delete(context.Context, *entity.PostID) error
}
