package domain

import (
	"time"
)

type Post struct {
	postID     *PostID
	postBody   *PostBody
	createTime time.Time
	updateTime time.Time
	deleteTime time.Time
}

func NewPost(postID *PostID, postBody *PostBody, createTime time.Time, updateTime time.Time, deleteTime time.Time) *Post {
	return &Post{
		postID:     postID,
		postBody:   postBody,
		createTime: createTime,
		updateTime: updateTime,
		deleteTime: deleteTime,
	}
}

func (p *Post) PostID() *PostID {
	return p.postID
}

func (p *Post) PostBody() *PostBody {
	return p.postBody
}

func (p *Post) CreateTime() time.Time {
	return p.createTime
}

func (p *Post) UpdateTime() time.Time {
	return p.updateTime
}

func (p *Post) DeleteTime() time.Time {
	return p.deleteTime
}

func (p *Post) Valid() bool {
	return p.deleteTime.IsZero()
}
