package domain

type PostBody string

func NewPostBody(postBody string) *PostBody {
	b := PostBody(postBody)
	return &b
}

func (b *PostBody) String() string {
	if b == nil {
		return ""
	}
	return string(*b)
}
