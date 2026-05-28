package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

//go:generate moq -out paginator_mock.go -pkg interfaces . Paginator

// Paginator はサーバ側で保持するページネーション状態（Pagination）の
// 永続化操作をまとめたインターフェイス。
//
// PostQuerier / PostCommander とは別概念のため独立した interface に分けている。
// 実装は同じ PostgreSQL を使っているが、将来的に Pagination だけ
// Redis や inmemory に切り出すといった選択肢を残しておく狙いもある。
//
// クライアントに見せる next_page_token は PaginationID（UUID 文字列）であり、
// 内部のカーソル形式は PostCursor として Pagination.cursor に JSON で
// 詰めて保管する。この interface は「保存と取得」だけに責務を絞り、
// カーソル形式の知識は持たない。
type Paginator interface {
	GetPagination(context.Context, domain.PaginationID) (*domain.Pagination, error)
	SavePagination(context.Context, *domain.Pagination) error
}
