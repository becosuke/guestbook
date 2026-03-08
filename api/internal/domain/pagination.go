package domain

import "github.com/google/uuid"

type PaginationID uuid.UUID

func NewPaginationID(paginationID string) PaginationID {
	if paginationID == "" {
		return PaginationID(uuid.Nil)
	}
	return PaginationID(uuid.MustParse(paginationID))
}

func (p PaginationID) String() string {
	return uuid.UUID(p).String()
}

func (p PaginationID) IsZero() bool {
	return uuid.UUID(p) == uuid.Nil
}

type Pagination struct {
	paginationID PaginationID
	cursor       []byte
}

func NewPagination(paginationID PaginationID, cursor []byte) *Pagination {
	return &Pagination{
		paginationID: paginationID,
		cursor:       cursor,
	}
}

func (p *Pagination) PaginationID() PaginationID {
	return p.paginationID
}

func (p *Pagination) Cursor() []byte {
	return p.cursor
}
