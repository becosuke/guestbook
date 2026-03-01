//go:generate moq -out querier_mock.go -pkg interfaces . Querier
package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain/entity"
)

type Querier interface {
	Get(context.Context, *entity.PostID) (*entity.Post, error)
	Range(context.Context, *entity.PageOption) ([]*entity.Post, error)
}
