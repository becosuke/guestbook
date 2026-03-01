package post

type Post struct {
	postID *PostID
	body   *Body
}

func NewPost(postID *PostID, body *Body) *Post {
	return &Post{
		postID: postID,
		body:   body,
	}
}

func (p *Post) PostID() *PostID {
	return p.postID
}

func (p *Post) Body() *Body {
	return p.body
}
