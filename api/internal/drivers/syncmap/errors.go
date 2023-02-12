package syncmap

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound        = errors.New("not exists")
	ErrInvalidData     = errors.New("returns invalid data")
	ErrInvalidArgument = errors.New("invalid argument")
)
