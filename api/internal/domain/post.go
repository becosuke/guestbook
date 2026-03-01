package domain

type Post struct {
	postID   *PostID
	postBody *PostBody
}

func NewPost(postID *PostID, postBody *PostBody) *Post {
	return &Post{
		postID:   postID,
		postBody: postBody,
	}
}

func (p *Post) PostID() *PostID {
	return p.postID
}

func (p *Post) PostBody() *PostBody {
	return p.postBody
}
