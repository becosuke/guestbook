//go:generate moq -out querier_mock.go -pkg repository . Querier
package repository

import (
	"context"

	entity "github.com/becosuke/guestbook/api/internal/domain/entity/post"
)

type Querier interface {
	Get(context.Context, *entity.PostID) (*entity.Post, error)
	Range(context.Context, *entity.PageOption) ([]*entity.Post, error)
}
