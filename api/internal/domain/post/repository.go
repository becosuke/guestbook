package post

import (
	"context"
)

type Repository interface {
	Get(context.Context, *Serial) (*Post, error)
	Range(context.Context, *PageOption) ([]*Post, error)
	Create(context.Context, *Post) (*Serial, error)
	Update(context.Context, *Post) error
	Delete(context.Context, *Serial) error
}
