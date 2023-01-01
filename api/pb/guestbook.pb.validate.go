// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: guestbook.proto

package pb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GetPostRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetPostRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetPostRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetPostRequestMultiError,
// or nil if none found.
func (m *GetPostRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetPostRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Serial

	if len(errors) > 0 {
		return GetPostRequestMultiError(errors)
	}

	return nil
}

// GetPostRequestMultiError is an error wrapping multiple validation errors
// returned by GetPostRequest.ValidateAll() if the designated constraints
// aren't met.
type GetPostRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetPostRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetPostRequestMultiError) AllErrors() []error { return m }

// GetPostRequestValidationError is the validation error returned by
// GetPostRequest.Validate if the designated constraints aren't met.
type GetPostRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetPostRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetPostRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetPostRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetPostRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetPostRequestValidationError) ErrorName() string { return "GetPostRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetPostRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetPostRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetPostRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetPostRequestValidationError{}

// Validate checks the field values on CreatePostRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CreatePostRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreatePostRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreatePostRequestMultiError, or nil if none found.
func (m *CreatePostRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePostRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPost() == nil {
		err := CreatePostRequestValidationError{
			field:  "Post",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetPost()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreatePostRequestValidationError{
					field:  "Post",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreatePostRequestValidationError{
					field:  "Post",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPost()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreatePostRequestValidationError{
				field:  "Post",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CreatePostRequestMultiError(errors)
	}

	return nil
}

// CreatePostRequestMultiError is an error wrapping multiple validation errors
// returned by CreatePostRequest.ValidateAll() if the designated constraints
// aren't met.
type CreatePostRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePostRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePostRequestMultiError) AllErrors() []error { return m }

// CreatePostRequestValidationError is the validation error returned by
// CreatePostRequest.Validate if the designated constraints aren't met.
type CreatePostRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePostRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePostRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePostRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePostRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePostRequestValidationError) ErrorName() string {
	return "CreatePostRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePostRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePostRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePostRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePostRequestValidationError{}

// Validate checks the field values on UpdatePostRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UpdatePostRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdatePostRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdatePostRequestMultiError, or nil if none found.
func (m *UpdatePostRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdatePostRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPost() == nil {
		err := UpdatePostRequestValidationError{
			field:  "Post",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetPost()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdatePostRequestValidationError{
					field:  "Post",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdatePostRequestValidationError{
					field:  "Post",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPost()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdatePostRequestValidationError{
				field:  "Post",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetFieldMask()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdatePostRequestValidationError{
					field:  "FieldMask",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdatePostRequestValidationError{
					field:  "FieldMask",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetFieldMask()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdatePostRequestValidationError{
				field:  "FieldMask",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UpdatePostRequestMultiError(errors)
	}

	return nil
}

// UpdatePostRequestMultiError is an error wrapping multiple validation errors
// returned by UpdatePostRequest.ValidateAll() if the designated constraints
// aren't met.
type UpdatePostRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdatePostRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdatePostRequestMultiError) AllErrors() []error { return m }

// UpdatePostRequestValidationError is the validation error returned by
// UpdatePostRequest.Validate if the designated constraints aren't met.
type UpdatePostRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdatePostRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdatePostRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdatePostRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdatePostRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdatePostRequestValidationError) ErrorName() string {
	return "UpdatePostRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdatePostRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdatePostRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdatePostRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdatePostRequestValidationError{}

// Validate checks the field values on DeletePostRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DeletePostRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeletePostRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeletePostRequestMultiError, or nil if none found.
func (m *DeletePostRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeletePostRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Serial

	if len(errors) > 0 {
		return DeletePostRequestMultiError(errors)
	}

	return nil
}

// DeletePostRequestMultiError is an error wrapping multiple validation errors
// returned by DeletePostRequest.ValidateAll() if the designated constraints
// aren't met.
type DeletePostRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeletePostRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeletePostRequestMultiError) AllErrors() []error { return m }

// DeletePostRequestValidationError is the validation error returned by
// DeletePostRequest.Validate if the designated constraints aren't met.
type DeletePostRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeletePostRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeletePostRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeletePostRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeletePostRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeletePostRequestValidationError) ErrorName() string {
	return "DeletePostRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeletePostRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeletePostRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeletePostRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeletePostRequestValidationError{}

// Validate checks the field values on ListPostsRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ListPostsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListPostsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListPostsRequestMultiError, or nil if none found.
func (m *ListPostsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListPostsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPageSize() <= 0 {
		err := ListPostsRequestValidationError{
			field:  "PageSize",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetPageToken() != "" {

	}

	if len(errors) > 0 {
		return ListPostsRequestMultiError(errors)
	}

	return nil
}

// ListPostsRequestMultiError is an error wrapping multiple validation errors
// returned by ListPostsRequest.ValidateAll() if the designated constraints
// aren't met.
type ListPostsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListPostsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListPostsRequestMultiError) AllErrors() []error { return m }

// ListPostsRequestValidationError is the validation error returned by
// ListPostsRequest.Validate if the designated constraints aren't met.
type ListPostsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListPostsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListPostsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListPostsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListPostsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListPostsRequestValidationError) ErrorName() string { return "ListPostsRequestValidationError" }

// Error satisfies the builtin error interface
func (e ListPostsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListPostsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListPostsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListPostsRequestValidationError{}

// Validate checks the field values on ListPostsResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ListPostsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListPostsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListPostsResponseMultiError, or nil if none found.
func (m *ListPostsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListPostsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPosts() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListPostsResponseValidationError{
						field:  fmt.Sprintf("Posts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListPostsResponseValidationError{
						field:  fmt.Sprintf("Posts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListPostsResponseValidationError{
					field:  fmt.Sprintf("Posts[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for NextPageToken

	if len(errors) > 0 {
		return ListPostsResponseMultiError(errors)
	}

	return nil
}

// ListPostsResponseMultiError is an error wrapping multiple validation errors
// returned by ListPostsResponse.ValidateAll() if the designated constraints
// aren't met.
type ListPostsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListPostsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListPostsResponseMultiError) AllErrors() []error { return m }

// ListPostsResponseValidationError is the validation error returned by
// ListPostsResponse.Validate if the designated constraints aren't met.
type ListPostsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListPostsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListPostsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListPostsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListPostsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListPostsResponseValidationError) ErrorName() string {
	return "ListPostsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListPostsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListPostsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListPostsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListPostsResponseValidationError{}

// Validate checks the field values on Post with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Post) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Post with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in PostMultiError, or nil if none found.
func (m *Post) ValidateAll() error {
	return m.validate(true)
}

func (m *Post) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Serial

	if l := utf8.RuneCountInString(m.GetBody()); l < 1 || l > 128 {
		err := PostValidationError{
			field:  "Body",
			reason: "value length must be between 1 and 128 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return PostMultiError(errors)
	}

	return nil
}

// PostMultiError is an error wrapping multiple validation errors returned by
// Post.ValidateAll() if the designated constraints aren't met.
type PostMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PostMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PostMultiError) AllErrors() []error { return m }

// PostValidationError is the validation error returned by Post.Validate if the
// designated constraints aren't met.
type PostValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PostValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PostValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PostValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PostValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PostValidationError) ErrorName() string { return "PostValidationError" }

// Error satisfies the builtin error interface
func (e PostValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPost.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PostValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PostValidationError{}
