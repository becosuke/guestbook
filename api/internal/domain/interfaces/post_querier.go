package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

//go:generate moq -out post_querier_mock.go -pkg interfaces . PostQuerier

// PostQuerier は Post の読み取り操作を表す CQS の Query 側インターフェイス。
//
// 書き込みを担う PostCommander と分離しているのは、
//   - 将来的に read replica からの参照に切り替える余地を残すため
//   - インターフェイス分離の原則（ISP）に沿って、読み取りだけ必要な
//     呼び出し側が書き込みメソッドに依存しないようにするため
//   - テスト時に「読み取りだけ」「書き込みだけ」をモック差し替えしやすくするため
// という意図。
type PostQuerier interface {
	// GetPost は PostID で指定された投稿を 1 件取得する。
	// 該当が無ければ domain.ErrNotFound を返す契約。
	GetPost(context.Context, domain.PostID) (*domain.Post, error)
	// RangePosts はカーソルベースで投稿を一覧取得する。
	// pageSize は「次ページの有無を判定するため +1 件多めに取る」運用なので、
	// 呼び出し側は要求件数 + 1 を渡すこと（実装の責務分担は usecase 側に置く）。
	RangePosts(ctx context.Context, pageSize int32, cursor *domain.PostCursor) ([]*domain.Post, error)
}
