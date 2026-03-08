package domain

import "github.com/google/uuid"

type PostID uuid.UUID

func NewPostID(postID string) PostID {
	if postID == "" {
		return PostID(uuid.Nil)
	}
	return PostID(uuid.MustParse(postID))
}

func (p PostID) String() string {
	return uuid.UUID(p).String()
}
