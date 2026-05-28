package interfaces

import (
	"context"

	"github.com/becosuke/guestbook/api/internal/domain"
)

//go:generate moq -out post_commander_mock.go -pkg interfaces . PostCommander

// PostCommander は Post の書き込み操作を表す CQS の Command 側インターフェイス。
//
// 読み取り（PostQuerier）と切り離されているのは、
//   - 書き込み経路だけマスタ DB に向ける構成を取りやすい
//   - 「副作用を持つ操作」だけを集めることでテスト時の影響範囲が明確になる
// から。
//
// 各メソッドは戻り値として Post を返さない。これは「書き込み後の最新状態を
// 取り直したい」という呼び出しは usecase 層が PostQuerier.GetPost を
// 続けて呼ぶ形で表現する設計にしているため（usecase.Create / Update を参照）。
type PostCommander interface {
	// CreatePost は新規 Post を永続化する。
	// 一意性制約に違反した場合は domain.ErrAlreadyExists を返す契約。
	CreatePost(context.Context, *domain.Post) error
	// UpdatePost は既存 Post の本文を更新し、旧本文を previousBody へ退避する。
	// 対象が存在しないか論理削除済みなら domain.ErrNotFound / ErrFailedPrecondition を返す。
	UpdatePost(context.Context, *domain.Post) error
	// DeletePost は論理削除を行う（deleteTime に現在時刻を入れる）。
	DeletePost(context.Context, domain.PostID) error
}
