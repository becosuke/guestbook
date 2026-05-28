package domain

import "strings"

// PostBody は投稿の本文を表す値オブジェクト。
// primitive obsession を避けるため、生 string ではなく薄いラッパ型として扱う。
type PostBody string

// NewPostBody は本文文字列から PostBody を構築する唯一のコンストラクタ。
//
// 前後の空白を TrimSpace で除去する正規化のみ行っている。これは
// 「本文末尾に改行が混ざっただけで等価性が崩れる」事態を防ぐためで、
// 長さや文字種といったバリデーションは presentation 層の protovalidate に任せる。
func NewPostBody(postBody string) PostBody {
	return PostBody(strings.TrimSpace(postBody))
}

// String は PostBody を生の文字列に戻す。
// 主にレスポンス組み立てと DB バインドで使われる。
func (b PostBody) String() string {
	return string(b)
}
