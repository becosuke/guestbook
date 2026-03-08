package interceptor

import (
	"context"
	"errors"

	"buf.build/go/protovalidate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ProtovalidateUnaryServerInterceptor returns a unary server interceptor that
// validates incoming requests using protovalidate.
func ProtovalidateUnaryServerInterceptor(validator protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if msg, ok := req.(proto.Message); ok {
			if err := validator.Validate(msg); err != nil {
				return nil, toGRPCError(err)
			}
		}
		return handler(ctx, req)
	}
}

// ProtovalidateStreamServerInterceptor returns a stream server interceptor that
// validates incoming messages using protovalidate.
func ProtovalidateStreamServerInterceptor(validator protovalidate.Validator) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		return handler(srv, &validatingServerStream{
			ServerStream: ss,
			validator:    validator,
		})
	}
}

type validatingServerStream struct {
	grpc.ServerStream
	validator protovalidate.Validator
}

func (s *validatingServerStream) RecvMsg(m any) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if msg, ok := m.(proto.Message); ok {
		if err := s.validator.Validate(msg); err != nil {
			return toGRPCError(err)
		}
	}
	return nil
}

func toGRPCError(err error) error {
	var valErr *protovalidate.ValidationError
	if !errors.As(err, &valErr) {
		return status.Error(codes.Internal, err.Error())
	}

	st := status.New(codes.InvalidArgument, valErr.Error())

	br := &errdetails.BadRequest{}
	for _, v := range valErr.Violations {
		br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       protovalidate.FieldPathString(v.Proto.GetField()),
			Description: v.Proto.GetMessage(),
		})
	}

	stWithDetails, err := st.WithDetails(br)
	if err != nil {
		return st.Err()
	}
	return stWithDetails.Err()
}
