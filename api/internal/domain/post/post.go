package post

type Post struct {
	serial *Serial
	body   *Body
}

func NewPost(serial *Serial, body *Body) *Post {
	return &Post{
		serial: serial,
		body:   body,
	}
}

func (p *Post) Serial() *Serial {
	return p.serial
}

func (p *Post) Body() *Body {
	return p.body
}
