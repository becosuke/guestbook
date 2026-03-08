package interceptor

import (
	"context"
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

func newValidator(t *testing.T) protovalidate.Validator {
	t.Helper()
	v, err := protovalidate.New()
	require.NoError(t, err)
	return v
}

func TestProtovalidateUnaryServerInterceptor_ValidRequest(t *testing.T) {
	v := newValidator(t)
	interceptor := ProtovalidateUnaryServerInterceptor(v)

	handlerCalled := false
	handler := func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return "ok", nil
	}

	req := &pb.GetPostRequest{PostId: "550e8400-e29b-41d4-a716-446655440000"}
	resp, err := interceptor(context.Background(), req, &grpc.UnaryServerInfo{}, handler)

	require.NoError(t, err)
	assert.Equal(t, "ok", resp)
	assert.True(t, handlerCalled)

	// CreatePostRequest with nil UUID should also be valid
	handlerCalled = false
	createReq := &pb.CreatePostRequest{
		Post:           &pb.Post{PostId: uuid.Nil.String(), Body: "hello"},
		IdempotencyKey: testUUID,
	}
	resp, err = interceptor(context.Background(), createReq, &grpc.UnaryServerInfo{}, handler)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp)
	assert.True(t, handlerCalled)
}

func TestProtovalidateUnaryServerInterceptor_InvalidRequest(t *testing.T) {
	v := newValidator(t)
	interceptor := ProtovalidateUnaryServerInterceptor(v)

	handlerCalled := false
	handler := func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return "ok", nil
	}

	tests := []struct {
		name    string
		req     any
		wantErr bool
	}{
		{
			name:    "invalid UUID in GetPostRequest",
			req:     &pb.GetPostRequest{PostId: "not-a-uuid"},
			wantErr: true,
		},
		{
			name:    "empty post_id in GetPostRequest",
			req:     &pb.GetPostRequest{PostId: ""},
			wantErr: true,
		},
		{
			name: "non-nil UUID post_id in CreatePostRequest",
			req: &pb.CreatePostRequest{
				Post:           &pb.Post{PostId: testUUID, Body: "hello"},
				IdempotencyKey: testUUID,
			},
			wantErr: true,
		},
		{
			name: "empty body in CreatePostRequest",
			req: &pb.CreatePostRequest{
				Post:           &pb.Post{PostId: uuid.Nil.String(), Body: ""},
				IdempotencyKey: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: true,
		},
		{
			name: "missing post in CreatePostRequest",
			req: &pb.CreatePostRequest{
				IdempotencyKey: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: true,
		},
		{
			name: "invalid idempotency_key in CreatePostRequest",
			req: &pb.CreatePostRequest{
				Post:           &pb.Post{Body: "hello"},
				IdempotencyKey: "not-a-uuid",
			},
			wantErr: true,
		},
		{
			name:    "zero page_size in ListPostsRequest",
			req:     &pb.ListPostsRequest{PageSize: 0},
			wantErr: true,
		},
		{
			name:    "negative page_size in ListPostsRequest",
			req:     &pb.ListPostsRequest{PageSize: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled = false
			resp, err := interceptor(context.Background(), tt.req, &grpc.UnaryServerInfo{}, handler)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, resp)
				assert.False(t, handlerCalled, "handler should not be called on validation error")

				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())

				details := st.Details()
				require.NotEmpty(t, details, "error should contain BadRequest details")
				br, ok := details[0].(*errdetails.BadRequest)
				require.True(t, ok)
				assert.NotEmpty(t, br.GetFieldViolations())
			} else {
				require.NoError(t, err)
				assert.True(t, handlerCalled)
			}
		})
	}
}

func TestProtovalidateUnaryServerInterceptor_BodyLength(t *testing.T) {
	v := newValidator(t)
	interceptor := ProtovalidateUnaryServerInterceptor(v)

	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "min length body (1 char)",
			body:    "a",
			wantErr: false,
		},
		{
			name:    "max length body (128 chars)",
			body:    strings.Repeat("a", 128),
			wantErr: false,
		},
		{
			name:    "over max length body (129 chars)",
			body:    strings.Repeat("a", 129),
			wantErr: true,
		},
		{
			name:    "empty body",
			body:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.CreatePostRequest{
				Post:           &pb.Post{PostId: uuid.Nil.String(), Body: tt.body},
				IdempotencyKey: testUUID,
			}
			resp, err := interceptor(context.Background(), req, &grpc.UnaryServerInfo{}, handler)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, resp)
				st, _ := status.FromError(err)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			} else {
				require.NoError(t, err)
				assert.Equal(t, "ok", resp)
			}
		})
	}
}

func TestToGRPCError_FieldViolationDetails(t *testing.T) {
	v := newValidator(t)

	req := &pb.CreatePostRequest{
		Post:           &pb.Post{Body: ""},
		IdempotencyKey: "not-a-uuid",
	}
	err := v.Validate(req)
	require.Error(t, err)

	grpcErr := toGRPCError(err)
	st, ok := status.FromError(grpcErr)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	details := st.Details()
	require.Len(t, details, 1)
	br, ok := details[0].(*errdetails.BadRequest)
	require.True(t, ok)
	assert.GreaterOrEqual(t, len(br.GetFieldViolations()), 2, "should have violations for both body and idempotency_key")

	fields := make(map[string]bool)
	for _, fv := range br.GetFieldViolations() {
		fields[fv.GetField()] = true
		assert.NotEmpty(t, fv.GetDescription())
	}
	assert.True(t, fields["post.body"], "should have violation for post.body")
	assert.True(t, fields["idempotency_key"], "should have violation for idempotency_key")
}
