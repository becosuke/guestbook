package syncmap

import (
	"github.com/pkg/errors"
)

var (
	ErrSyncmapNotFound        = errors.New("not exists")
	ErrSyncmapInvalidData     = errors.New("returns invalid data")
	ErrSyncmapInvalidArgument = errors.New("invalid argument")
)
