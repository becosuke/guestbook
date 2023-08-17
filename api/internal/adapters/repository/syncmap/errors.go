package syncmap

import (
	"github.com/pkg/errors"
)

var (
	ErrMessageAlreadyExists = errors.New("message already exists")
	ErrMessageNotFound      = errors.New("not exists")
	ErrInvalidData          = errors.New("returns invalid data")
	ErrInvalidArgument      = errors.New("invalid argument")
)
