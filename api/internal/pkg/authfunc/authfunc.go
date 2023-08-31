package authfunc

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

func NewAuthFunc(ctx context.Context) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	}
}
