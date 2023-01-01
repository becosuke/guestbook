package syncmap

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Errors interface {
	error
	GRPCStatus() *status.Status
}

type errorImpl struct {
	codes.Code
	msg string
}

func NewSyncmapError(code codes.Code, msg string) Errors {
	return &errorImpl{Code: code, msg: msg}
}

func (impl *errorImpl) GRPCStatus() *status.Status {
	return status.New(impl.Code, impl.msg)
}

func (impl *errorImpl) Error() string {
	return impl.String()
}

var (
	ErrSyncmapNotFound        = NewSyncmapError(codes.NotFound, "not exists")
	ErrSyncmapInvalidData     = NewSyncmapError(codes.Internal, "returns invalid data")
	ErrSyncmapInvalidArgument = NewSyncmapError(codes.InvalidArgument, "invalid argument")
)
