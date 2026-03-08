package domain

import "strings"

type PostBody string

func NewPostBody(postBody string) PostBody {
	return PostBody(strings.TrimSpace(postBody))
}

func (b PostBody) String() string {
	return string(b)
}
