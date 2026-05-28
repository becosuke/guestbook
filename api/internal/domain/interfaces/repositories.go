// Package interfaces はドメイン層が外部世界に依存するために定義する
// リポジトリインターフェイス群を置く場所。
//
// 依存性逆転の原則（DIP）に従い、ドメイン／usecase 層は「自分が必要とする
// 操作」を interface で宣言し、その実装は adapter/repository 側が満たす。
// これによりドメイン層は PostgreSQL や Redis といった具体技術を一切知らずに済む。
//
// 各インターフェイスは moq を用いてモックを自動生成しており、テスト時は
// 関数フィールドを差し替えるだけで任意の挙動を再現できる。
package interfaces

//go:generate moq -out repositories_mock.go -pkg interfaces . Repositories

// Repositories は usecase 層が依存する集約インターフェイス。
//
// 個々の Querier / Commander / Paginator を usecase に直接注入する形でも
// 機能上は十分だが、
//   - 依存注入のコンストラクタ引数を 1 つにまとめられる
//   - 「永続化の総体」という単一の抽象として扱える
// 利点を取って 1 つの interface に embed している。
//
// CQS の観点で読み取り（Querier）と書き込み（Commander）を別 interface に
// 分けたままにしているので、テストでは Repositories ではなく個別の
// 細粒度モックを使うことも可能。
type Repositories interface {
	PostQuerier
	PostCommander
	Paginator
}
