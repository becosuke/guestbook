package domain

import (
	"encoding/json"
	"time"
)

type PostCursor struct {
	LastPostID     string    `json:"last_post_id"`
	LastCreateTime time.Time `json:"last_create_time"`
}

func NewPostCursor(lastPostID string, lastCreateTime time.Time) *PostCursor {
	return &PostCursor{
		LastPostID:     lastPostID,
		LastCreateTime: lastCreateTime,
	}
}

func (c *PostCursor) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func UnmarshalPostCursor(data []byte) (*PostCursor, error) {
	var cursor PostCursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return nil, err
	}
	return &cursor, nil
}
