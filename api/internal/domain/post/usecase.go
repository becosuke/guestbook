package post

import (
	"context"
)

type Usecase interface {
	Get(context.Context, *Serial) (*Post, error)
	Range(context.Context, *PageOption) ([]*Post, error)
	Create(context.Context, *Post) (*Post, error)
	Update(context.Context, *Post) (*Post, error)
	Delete(context.Context, *Serial) error
}
