package post

type PageSize int32
type PageToken string

type PageOption struct {
	pageSize  *PageSize
	pageToken *PageToken
}

func NewPageOption(pageSize *PageSize, pageToken *PageToken) *PageOption {
	return &PageOption{
		pageSize:  pageSize,
		pageToken: pageToken,
	}
}

func (p *PageOption) PageSize() *PageSize {
	return p.pageSize
}

func (p *PageOption) PageToken() *PageToken {
	return p.pageToken
}
