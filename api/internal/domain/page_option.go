package domain

// このファイル全体は Google AIP-158（Pagination）の規約を素直に
// ドメイン層へ落とし込んだもの。AIP は List 系メソッドの page_size /
// page_token の振る舞いを標準化しており、それに従うことでクライアント
// 実装やドキュメントの予測可能性が高まる。

// PageSize はクライアントが要求する 1 ページあたりの件数を表す。
//
// int32 を採用しているのは AIP-158 が page_size を int32 と規定しているため。
// 「int64 の方が表現範囲が広くて安全では」と思える局面でも、ここは仕様に
// 合わせる方を優先する。protobuf 型と一致するためバイナリ表現上の取り扱いも
// 自然になる。
type PageSize int32

// PageToken は次ページを取得するための不透明（opaque）なトークン。
//
// AIP-158 では「クライアントはトークンの中身を解釈してはならず、サーバから
// 受け取った値をそのまま次のリクエストに渡す」と定めている。実体は
// PaginationID の文字列表現だが、それはサーバ側の実装詳細で、
// クライアントから見ると「不透明な文字列」として扱われる。
type PageToken string

// PageOption は List 系メソッドのページネーション指定をまとめた値。
//
// AIP-158 のセマンティクスに沿った設計上のポイント:
//   - page_size の解釈は「クライアントの希望値」であり、サーバは
//     最終的な件数（デフォルトや上限）を自分で決めて構わない
//   - page_token が空ならば「最初のページ」を意味する
//   - レスポンスの next_page_token が空ならば「これ以上ページが無い」を意味する
//
// 各フィールドをポインタにしているのは「未指定（nil）」と「明示的に 0 や
// 空文字を指定された」を区別するため。
// 実運用では converter が常に非 nil で値を詰めるため、nil 経路はほぼ通らないが、
// page_size に関しては protovalidate を int32.gte=0 にしてあるので 0 が
// usecase まで素通りする。usecase 側で「nil または 0 ならサーバデフォルト
// （page_size=10）を適用」と扱う契約になっている。
type PageOption struct {
	pageSize  *PageSize
	pageToken *PageToken
}

// NewPageOption は PageOption の唯一のコンストラクタ。
// 引数の nil 許容は上記コメントの通り意図的。
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
