package repository

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrInvalidData     = errors.New("returns invalid data")
)
