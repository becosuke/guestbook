package domain

import "github.com/google/uuid"

// このファイルは Google AIP-158（Pagination）が要求する page_token の
// 性質を、サーバ保持型カーソルとして実現するためのドメインモデル。
//
// AIP-158 はクライアントから見て page_token を「中身を解釈してはならない
// 不透明（opaque）な文字列」と定めている。これを満たす実装方式には
//   (a) カーソル情報を暗号化／署名してトークンに埋め込む方式
//   (b) サーバ側にカーソル情報を保存し、トークンにはその識別子だけを返す方式
// があり、本プロジェクトでは (b) を採用した。決め手は次の 3 点:
//   1. クライアントから渡されるトークンの改ざんリスクをそもそも遮断できる
//   2. トークン長が UUID 文字列で一定（埋め込み方式は内容に応じて伸縮する）
//   3. 内部カーソル形式の変更がクライアント互換に影響しない
//
// page_token の空文字を「最初のページ」、next_page_token の空文字を
// 「終端」とする AIP-158 のセマンティクスは、PaginationID のゼロ値で
// 表現している（IsZero / NewPaginationID("") を参照）。

// PaginationID はサーバ側で保持しているページネーション状態の識別子。
// 設計思想は PostID と同じく typed UUID。
type PaginationID uuid.UUID

// NewPaginationID は文字列から PaginationID を構築する唯一のコンストラクタ。
// 不正な UUID 文字列に対しては NewPostID と同じく MustParse で panic する
// 方針を取り、ドメインに重複したバリデーションを置かない。
//
// 空文字を受け取ったときは uuid.Nil を返し、これが PaginationID のゼロ表現になる。
// 空文字をそのまま保持せず uuid.Nil に揃えているのは、
//   1. typed UUID として「該当する pagination が存在しない」ことを明示的に
//      表現するため
//   2. uuid パッケージが UUID のゼロ値を uuid.Nil と定めており、Go における
//      `var x PaginationID` のゼロ値と整合させるため
// の 2 点による。外側から見ると「page_token が空文字 = 最初のページ」という
// AIP-158 のセマンティクスとも対応する。
func NewPaginationID(paginationID string) PaginationID {
	if paginationID == "" {
		return PaginationID(uuid.Nil)
	}
	return PaginationID(uuid.MustParse(paginationID))
}

func (p PaginationID) String() string {
	return uuid.UUID(p).String()
}

// IsZero は「これ以上ページが続かない」状態を判定する。
//
// 実体は「PaginationID が uuid.Nil（= 型のゼロ値）かどうか」を見ているだけ。
// usecase.Range は最終ページでこのゼロ値を返し、presentation 層は IsZero の
// 結果で next_page_token を出力するかを判断する。AIP-158 が定める
// 「next_page_token が空文字なら終端」というセマンティクスに対応する。
func (p PaginationID) IsZero() bool {
	return uuid.UUID(p) == uuid.Nil
}

// Pagination はサーバ側で保管するページネーション状態のエンティティ。
//
// PaginationID をキーに cursor（次ページの開始位置情報。実体は PostCursor を
// JSON シリアライズしたバイト列）を紐付けて永続化することで、クライアントには
// 軽量な ID だけを返し、内部状態はサーバが管理する。これがファイル冒頭で
// 述べたサーバ保持型カーソルの本体。
//
// cursor が []byte なのは、Pagination が「カーソルの中身を解釈しない」という
// 関心の分離を意図しているため。中身の構造は PostCursor 側に閉じる。
type Pagination struct {
	paginationID PaginationID
	cursor       []byte
}

// NewPagination は Pagination の唯一のコンストラクタ。
// usecase 層が新しいカーソルを発行するときと、repository 層が DB から
// 読み出した行を再構成するとき、両方の経路で使う。
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
