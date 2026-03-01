package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostID_String(t *testing.T) {
	type testCase struct {
		name   string
		postID *PostID
		want   string
	}
	tests := []testCase{
		func() testCase {
			return testCase{
				name:   "normal",
				postID: NewPostID("550e8400-e29b-41d4-a716-446655440000"),
				want:   "550e8400-e29b-41d4-a716-446655440000",
			}
		}(),
		func() testCase {
			return testCase{
				name:   "empty",
				postID: NewPostID(""),
				want:   "",
			}
		}(),
		func() testCase {
			return testCase{
				name:   "nil",
				postID: nil,
				want:   "",
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.postID.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostBody_String(t *testing.T) {
	type testCase struct {
		name     string
		postBody *PostBody
		want     string
	}
	tests := []testCase{
		func() testCase {
			return testCase{
				name:     "normal",
				postBody: NewPostBody("example"),
				want:     "example",
			}
		}(),
		func() testCase {
			return testCase{
				name:     "empty",
				postBody: NewPostBody(""),
				want:     "",
			}
		}(),
		func() testCase {
			return testCase{
				name:     "nil",
				postBody: nil,
				want:     "",
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.postBody.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewPost(t *testing.T) {
	type args struct {
		postID   *PostID
		postBody *PostBody
	}
	type testCase struct {
		name         string
		args         args
		wantPostID   *PostID
		wantPostBody *PostBody
	}
	tests := []testCase{
		func() testCase {
			postID := NewPostID("550e8400-e29b-41d4-a716-446655440000")
			postBody := NewPostBody("example")
			return testCase{
				name: "normal",
				args: args{
					postID:   postID,
					postBody: postBody,
				},
				wantPostID:   postID,
				wantPostBody: postBody,
			}
		}(),
		func() testCase {
			return testCase{
				name: "nil fields",
				args: args{
					postID:   nil,
					postBody: nil,
				},
				wantPostID:   nil,
				wantPostBody: nil,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPost(tt.args.postID, tt.args.postBody)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantPostID, got.PostID())
			assert.Equal(t, tt.wantPostBody, got.PostBody())
		})
	}
}

func TestNewPageOption(t *testing.T) {
	type args struct {
		pageSize  *PageSize
		pageToken *PageToken
	}
	type testCase struct {
		name          string
		args          args
		wantPageSize  *PageSize
		wantPageToken *PageToken
	}
	tests := []testCase{
		func() testCase {
			pageSize := PageSize(10)
			pageToken := PageToken("token123")
			return testCase{
				name: "normal",
				args: args{
					pageSize:  &pageSize,
					pageToken: &pageToken,
				},
				wantPageSize:  &pageSize,
				wantPageToken: &pageToken,
			}
		}(),
		func() testCase {
			pageSize := PageSize(20)
			return testCase{
				name: "nil page token",
				args: args{
					pageSize:  &pageSize,
					pageToken: nil,
				},
				wantPageSize:  &pageSize,
				wantPageToken: nil,
			}
		}(),
		func() testCase {
			return testCase{
				name: "nil fields",
				args: args{
					pageSize:  nil,
					pageToken: nil,
				},
				wantPageSize:  nil,
				wantPageToken: nil,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPageOption(tt.args.pageSize, tt.args.pageToken)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantPageSize, got.PageSize())
			assert.Equal(t, tt.wantPageToken, got.PageToken())
		})
	}
}
