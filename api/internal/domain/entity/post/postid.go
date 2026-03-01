package post

type PostID string

func NewPostID(postID string) *PostID {
	p := PostID(postID)
	return &p
}

func (p *PostID) String() string {
	if p == nil {
		return ""
	}
	return string(*p)
}
