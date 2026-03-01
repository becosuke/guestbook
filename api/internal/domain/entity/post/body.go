package post

type Body string

func NewBody(body string) *Body {
	b := Body(body)
	return &b
}

func (b *Body) String() string {
	if b == nil {
		return ""
	}
	return string(*b)
}
