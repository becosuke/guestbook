package domain

type PaginationID string

func NewPaginationID(paginationID string) *PaginationID {
	p := PaginationID(paginationID)
	return &p
}

func (p *PaginationID) String() string {
	if p == nil {
		return ""
	}
	return string(*p)
}

type Pagination struct {
	paginationID *PaginationID
	cursor       []byte
}

func NewPagination(paginationID *PaginationID, cursor []byte) *Pagination {
	return &Pagination{
		paginationID: paginationID,
		cursor:       cursor,
	}
}

func (p *Pagination) PaginationID() *PaginationID {
	return p.paginationID
}

func (p *Pagination) Cursor() []byte {
	return p.cursor
}
