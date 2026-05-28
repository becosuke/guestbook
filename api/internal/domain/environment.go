package domain

import (
	"errors"
	"strings"
)

// Environment は実行環境を表す列挙的な値オブジェクト。
//
// 環境変数 ENVIRONMENT から読み出した生 string をそのまま引き回すと、
// 比較箇所ごとに "development" などのリテラルが点在して typo しやすい。
// typed string + 定数で表現することで、コンパイラと IDE 補完に
// 安全性を担保させる。
type Environment string

const (
	// EnvUnknown は環境判定に失敗した際のフォールバック。
	// ゼロ値のように扱われることを想定して用意している。
	EnvUnknown     Environment = "unknown"
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
	EnvTest        Environment = "test"
)

var (
	// ErrEnvUnknown は未知の環境名が指定されたときに返すドメインエラー。
	// 起動時の envconfig 解決失敗を errors.Is で識別するためのもの。
	ErrEnvUnknown = errors.New("unknown environment")
)

// NewEnvironment は文字列から Environment を構築する。
//
// PostID 系のコンストラクタが panic する設計だったのに対し、こちらは
// エラーを返すスタイルにしてある。
//
// この使い分けは「終了するかどうか」ではなく「呼び出し元のハンドリング契約」の
// 違いによる。Go の慣習では panic は「想定外の不整合（= バグ）」を表すもので、
// 通常は呼び出し元が個別に拾うことを期待しない。一方 error は型シグネチャに
// 失敗の可能性を出すことで、呼び出し元に「ここは失敗しうるので明示的に
// 扱ってほしい」と伝える契約になる。
//
// NewEnvironment が想定する失敗は ENVIRONMENT=prodcution のような
// 「運用者の入力ミス」であり、これはバグではなく起動時に十分起こりうる
// 通常のエラー。だから error で返し、infrastructure/config から main まで
// 持ち上げて、原因を整形した起動失敗ログを出してから exit code を選ばせる。
// PostID の NewPostID は presentation 層の protovalidate を通過した後に
// 不正値が届く想定が無く、それでも届いたらドメイン側の前提が崩れている
// （= バグ）ので panic が妥当、という対比になっている。
//
// なお、黙って EnvUnknown を返してしまうと、サーバが本番扱いされないまま
// 起動して、ログレベル・認証・接続先などが本番想定で切り替わらない事故に
// つながる。そのため Unknown はゼロ値表現に留め、コンストラクタは必ず
// エラーを伴って返す設計にしている。
//
// 大文字小文字は ToLower で正規化しているので、"Development" や "PRODUCTION"
// のような表記揺れは受け入れる。
func NewEnvironment(s string) (Environment, error) {
	switch strings.ToLower(s) {
	case "development":
		return EnvDevelopment, nil
	case "production":
		return EnvProduction, nil
	case "test":
		return EnvTest, nil
	default:
		return EnvUnknown, ErrEnvUnknown
	}
}

func (e Environment) String() string {
	return string(e)
}
